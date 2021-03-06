package system

import (
	"app.modules/system/customerror"
	"app.modules/system/myfirestore"
	"app.modules/system/mylinebot"
	"app.modules/system/youtubebot"
	"context"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	InCommand = "!in"
	OutCommand = "!out"
	InfoCommand = "!info"
	CommandPrefix = "!"

	WorkNameOptionPrefix = "work-"
	WorkTimeOptionPrefix = "min-"

	FullWidthSpace = "　"
	HalfWidthSpace = " "

)

type System struct {
	FirestoreController *myfirestore.FirestoreController
	LiveChatBot *youtubebot.YoutubeLiveChatBot
	LineBot *mylinebot.LineBot
	MinWorkTimeMin int
	MaxWorkTimeMin int
	ProcessedUserId string
	ProcessedUserDisplayName string
	DefaultSleepIntervalMilli int
}

func NewSystem(ctx context.Context, clientOption option.ClientOption) (System, error) {
	fsController, err := myfirestore.NewFirestoreController(ctx, ProjectId, clientOption)
	if err != nil {
		return System{}, err
	}

	// youtube live chat bot
	youtubeLiveConfig, err := fsController.RetrieveYoutubeLiveConfig(ctx)
	if err != nil {
		return System{}, err
	}
	liveChatBot, err := youtubebot.NewYoutubeLiveChatBot(youtubeLiveConfig.LiveChatId, fsController, ctx)
	if err != nil {
		return System{}, err
	}

	// line bot
	lineBotConfig, err := fsController.RetrieveLineBotConfig(ctx)
	if err != nil {
		return System{}, err
	}
	lineBot, err := mylinebot.NewLineBot(lineBotConfig.ChannelSecret, lineBotConfig.ChannelToken, lineBotConfig.DestinationLineId)
	if err != nil {
		return System{}, err
	}

	// system constant values
	constantsConfig, err := fsController.RetrieveSystemConstantsConfig(ctx)
	if err != nil {
		return System{}, err
	}

	return System{
		FirestoreController: fsController,
		LiveChatBot:         liveChatBot,
		LineBot:             lineBot,
		MaxWorkTimeMin: constantsConfig.MaxWorkTimeMin,
		MinWorkTimeMin: constantsConfig.MinWorkTimeMin,
		DefaultSleepIntervalMilli: constantsConfig.SleepIntervalMilli,
	}, nil
}

func (s *System) SetProcessedUser(userId string, userDisplayName string) {
	s.ProcessedUserId = userId
	s.ProcessedUserDisplayName = userDisplayName
}

func (s *System) CloseFirestoreClient() {
	err := s.FirestoreController.FirestoreClient.Close()
	if err != nil {
		log.Println("failed close firestore client.")
	} else {
		log.Println("successfully closed firestore client.")
	}
}

// Command: 入力コマンドを解析して実行
func (s *System) Command(commandString string, userId string, userDisplayName string, ctx context.Context) error {
	if strings.HasPrefix(commandString, CommandPrefix) {
		s.SetProcessedUser(userId, userDisplayName)
		slice := strings.Split(commandString, HalfWidthSpace)
		switch slice[0] {
		case InCommand:
			return s.In(commandString, ctx)
		case OutCommand:
			return s.Out(ctx)
		case InfoCommand:
			return s.ShowUserInfo(ctx)
		default:
			// !席番号
			num, err := strconv.Atoi(strings.TrimLeft(slice[0], CommandPrefix))
			if err == nil && num >= 0 {
				return s.In(commandString, ctx)
			}
		}
	}
	return nil
}

func (s *System) In(commandString string, ctx context.Context) error {
	slice := strings.Split(commandString, HalfWidthSpace)

	// すでに入室している場合
	isInRoom, err := s.IsUserInRoom(ctx)
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed s.IsUserInRoom()", err)
		return err
	}
	if isInRoom {
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さん、すでに入室しています。まず「" + OutCommand + "」で退室してください。", ctx)
		return nil
	}

	// 席を指定しているかどうか
	var seatId int
	num, err := strconv.Atoi(strings.TrimLeft(slice[0], CommandPrefix))
	if err != nil {	// !in
		seatId, err = s.RandomAvailableSeatId(ctx)
		if err != nil {
			s.SendLiveChatMessage(s.ProcessedUserDisplayName +
				"さん、エラーが発生しました。もう一度試してみてください。", ctx)
			return err
		}
	} else {	// 指定された座席番号が有効かチェック
		seatId = num
		switch seatId {
		case 0:
			break
		default:
			// その席番号が存在するか
			isSeatExist, err := s.IsSeatExist(seatId, ctx)
			if err != nil {
				s.SendLiveChatMessage(s.ProcessedUserDisplayName +
					"さん、エラーが発生しました。もう一度試してみてください。", ctx)
				_ = s.LineBot.SendMessageWithError("failed s.IsSeatExist()", err)
				return err
			} else if ! isSeatExist {
				s.SendLiveChatMessage(s.ProcessedUserDisplayName + "さん、その番号の席は" +
					"存在しません。他の空いている席を選ぶか、「" + InfoCommand + "」で席を指定せずに入室してください！", ctx)
				return nil
			}
			// その席が空いているか
			isOk, err := s.IfSeatAvailable(seatId, ctx)
			if err != nil {
				_ = s.LineBot.SendMessageWithError("failed s.IfSeatAvailable()", err)
				s.SendLiveChatMessage(s.ProcessedUserDisplayName +
					"さん、エラーが発生しました。もう一度試してみてください。", ctx)
				return err
			}
			if ! isOk {
				s.SendLiveChatMessage(s.ProcessedUserDisplayName +
					"さん、その席には今は座れません！空いている座席の番号を書いてください！", ctx)
				return nil
			}
		}
	}

	// 追加オプションチェック
	workName := ""
	workTimeMin := 120
	for _, str := range slice[1:] {
		if strings.HasPrefix(str, WorkNameOptionPrefix) {
			workName = strings.TrimLeft(str, WorkNameOptionPrefix)
		} else if strings.HasPrefix(str, WorkTimeOptionPrefix) {
			num, err = strconv.Atoi(strings.TrimLeft(str, WorkTimeOptionPrefix))
			if err != nil {
				s.SendLiveChatMessage(s.ProcessedUserDisplayName +
					"さん、「" + WorkTimeOptionPrefix + "」の後の数字は半角になっているか確認してみてください。", ctx)
				return nil
			}
			if 5 <= num && num <= 360 {
				workTimeMin = num
			} else {
				s.SendLiveChatMessage(s.ProcessedUserDisplayName +
					"さん、作業時間（分）は" + strconv.Itoa(s.MinWorkTimeMin) + "～" + strconv.Itoa(s.MaxWorkTimeMin) + "の値にしてください。", ctx)
				return nil
			}
		}
	}
	
	// 初回の利用の場合はユーザーデータを初期化
	isRegistered, err := s.IfUserRegistered(ctx)
	if err != nil {
		return err
	}
	if ! isRegistered {
		err := s.InitializeUser(ctx)
		if err != nil {
			return err
		}
	}

	// 入室
	if seatId == 0 {	// no-seat-room
		err = s.EnterNoSeatRoom(workName, workTimeMin, ctx)
	} else { // default-room
		err = s.EnterDefaultRoom(seatId, workName, workTimeMin, ctx)
	}
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed to enter room", err)
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さん、エラーが発生しました。もう一度試してみてください。", ctx)
		return err
	} else {
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さんが作業を始めました！（最大" + strconv.Itoa(workTimeMin) + "分）", ctx)
		//s.SendLiveChatMessage(s.ProcessedUserDisplayName +
		//	" started working!! (" + strconv.Itoa(workTimeMin) + " minutes max.)", ctx)
		// 入室時刻を記録
		err = s.FirestoreController.SetLastEnteredDate(s.ProcessedUserId, ctx)
		if err != nil {
			_ = s.LineBot.SendMessageWithError("failed to set last entered date", err)
			s.SendLiveChatMessage(s.ProcessedUserDisplayName +
				"さん、エラーが発生しました。もう一度試してみてください。", ctx)
			return err
		}
	}
	return nil
}

func (s *System) Out(ctx context.Context) error {
	// 今勉強中か？
	isInRoom, err := s.IsUserInRoom(ctx)
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed IsUserInRoom()", err)
		return err
	}
	if ! isInRoom {
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さん、すでに退室してます！", ctx)
		return nil
	}
	// 退室処理
	err = s.ExitRoom(ctx)
	if err != nil {
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さん、エラーが発生しました。もう一度試してみてください。", ctx)
		return err
	}
	return nil
}

func (s *System) ShowUserInfo(ctx context.Context) error {
	// そのユーザーはデータがあるか？
	isUserRegistered, err := s.IfUserRegistered(ctx)
	if err != nil {
		return err
	}
	if isUserRegistered {
		totalTimeStr, dailyTotalTimeStr, err := s.TotalStudyTimeStrings(ctx)
		if err != nil {
			return err
		}
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さんの本日の作業時間は" + dailyTotalTimeStr + "、" +
			"累計作業時間は" + totalTimeStr + "です。", ctx)
		//s.SendLiveChatMessage("Hi, " + s.ProcessedUserDisplayName +
		//	". Your daily total working time is " + dailyTotalTimeStr + ", " +
		//	"cumulative working time is " + totalTimeStr + ".", ctx)
	} else {
		s.SendLiveChatMessage(s.ProcessedUserDisplayName +
			"さんはまだ作業データがありません。「" + InCommand + "」コマンドで作業を始めましょう！", ctx)
	}
	return nil
}

// IfSeatAvailable: 席番号がseatIdの席が空いているかどうか。seatIdは存在するという前提
func (s *System) IfSeatAvailable(seatId int, ctx context.Context) (bool, error) {
	defaultRoomData, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return false, err
	}
	for _, seat := range defaultRoomData.Seats {
		if seat.SeatId == seatId {
			return false, nil
		}
	}
	// ここまで来ると指定された番号の席が使われていないということ
	return true, nil
}

func (s *System) IsUserInRoom(ctx context.Context) (bool, error) {
	defaultRoomData, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return false, err
	}
	for _, seatInDefaultRoom := range defaultRoomData.Seats {
		if seatInDefaultRoom.UserId == s.ProcessedUserId  {
			return true, nil
		}
	}

	noSeatRoomData, err := s.FirestoreController.RetrieveNoSeatRoom(ctx)
	if err != nil {
		return false, err
	}
	for _, seatInNoSeatRoom := range noSeatRoomData.Seats {
		if seatInNoSeatRoom.UserId == s.ProcessedUserId  {
			return true, nil
		}
	}
	return false, nil
}

func (s *System) InitializeUser(ctx context.Context) error {
	log.Println("InitializeUser()")
	userData := myfirestore.UserDoc{
		DailyTotalStudySec: 0,
		TotalStudySec:      0,
		RegistrationDate:   time.Now(),
	}
	err := s.FirestoreController.InitializeUser(s.ProcessedUserId, userData, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *System) RetrieveYoutubeLiveInfo(ctx context.Context) (myfirestore.YoutubeLiveConfigDoc, error) {
	return s.FirestoreController.RetrieveYoutubeLiveConfig(ctx)
}

func (s *System) RetrieveNextPageToken(ctx context.Context) (string, error) {
	return s.FirestoreController.RetrieveNextPageToken(ctx)
}

func (s *System) SaveNextPageToken(nextPageToken string, ctx context.Context) error {
	return s.FirestoreController.SaveNextPageToken(nextPageToken, ctx)
}

func (s *System) EnterDefaultRoom(seatId int, workName string, workTimeMin int, ctx context.Context) error {
	exitDate := time.Now().Add(time.Duration(workTimeMin) * time.Minute)
	seat, err := s.FirestoreController.SetSeatInDefaultRoom(seatId, workName, exitDate, s.ProcessedUserId, s.ProcessedUserDisplayName, ctx)
	if err != nil {
		return err
	}
	// 入室時刻を記録
	err = s.FirestoreController.SetLastEnteredDate(s.ProcessedUserId, ctx)
	if err != nil {
		return err
	}
	// ログ記録
	err = s.FirestoreController.AddUserHistory(s.ProcessedUserId, EnterAction, seat, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *System) EnterNoSeatRoom(workName string, workTimeMin int, ctx context.Context) error {
	exitDate := time.Now().Add(time.Duration(workTimeMin) * time.Minute)
	seat, err := s.FirestoreController.SetSeatInNoSeatRoom(workName, exitDate, s.ProcessedUserId, s.ProcessedUserDisplayName, ctx)
	if err != nil {
		return err
	}
	// 退室時刻を記録
	err = s.FirestoreController.SetLastExitedDate(s.ProcessedUserId, ctx)
	if err != nil {
		return err
	}
	// ログ記録
	err = s.FirestoreController.AddUserHistory(s.ProcessedUserId, EnterAction, seat, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *System) RandomAvailableSeatId(ctx context.Context) (int, error) {
	roomLayout, err := s.FirestoreController.RetrieveDefaultRoomLayout(ctx)
	if err != nil {
		return 0, err
	}
	defaultRoom, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return 0, err
	}
	
	var availableSeatIdList []int
	for _, seatInLayout := range roomLayout.Seats {
		isUsed := false
		for _, seatInUse := range defaultRoom.Seats {
			if seatInLayout.Id == seatInUse.SeatId {
				isUsed = true
				break
			}
		}
		if ! isUsed {
			availableSeatIdList = append(availableSeatIdList, seatInLayout.Id)
		}
	}
	
	if len(availableSeatIdList) > 0 {
		rand.Seed(time.Now().UnixNano())
		return availableSeatIdList[rand.Intn(len(availableSeatIdList))], nil
	} else {
		return 0, nil
	}
}

func (s *System) ExitRoom(ctx context.Context) error {
	seatId, customErr := s.CurrentSeatId(ctx)
	if customErr.Body != nil {
		if customErr.ErrorType == customerror.UserNotInAnyRoom {
			s.SendLiveChatMessage(s.ProcessedUserDisplayName +
				"さん、あなたは今ルーム内にはいません。", ctx)
			return nil
		} else {
			s.SendLiveChatMessage(s.ProcessedUserDisplayName +
				"さん、残念ながらエラーが発生しました。もう一度試してみてください。", ctx)
			return customErr.Body
		}
	}
	// 作業時間を計算
	userData, err := s.FirestoreController.RetrieveUser(s.ProcessedUserId, ctx)
	if err != nil {
		return err
	}
	workedTimeSec := int(time.Now().Sub(userData.LastEntered).Seconds())

	var seat myfirestore.Seat
	switch seatId {
	case 0:
		noSeatRoom, err := s.FirestoreController.RetrieveNoSeatRoom(ctx)
		if err != nil {
			return err
		}
		for _, seatInNoSeatRoom := range noSeatRoom.Seats {
			if seatInNoSeatRoom.UserId == s.ProcessedUserId {
				seat = seatInNoSeatRoom
			}
		}
		err = s.FirestoreController.UnSetSeatInNoSeatRoom(seat, ctx)
		if err != nil {
			s.SendLiveChatMessage(s.ProcessedUserDisplayName +
				"さん、残念ながらエラーが発生しました。もう一度試してみてください。", ctx)
			return err
		}
		s.SendLiveChatMessage(s.ProcessedUserDisplayName + "さんが退室しました！" +
			"（作業時間" + strconv.Itoa(workedTimeSec / 60) + "分）", ctx)
		//s.SendLiveChatMessage(s.ProcessedUserDisplayName + " has finished working! " +
		//	"(" + strconv.Itoa(workedTimeSec / 60) + " minutes)", ctx)
	default:
		defaultSeatRoom, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
		if err != nil {
			return err
		}
		for _, seatDefaultRoom := range defaultSeatRoom.Seats {
			if seatDefaultRoom.UserId == s.ProcessedUserId {
				seat = seatDefaultRoom
			}
		}
		err = s.FirestoreController.UnSetSeatInDefaultRoom(seat, ctx)
		if err != nil {
			s.SendLiveChatMessage(s.ProcessedUserDisplayName +
				"さん、残念ながらエラーが発生しました。もう一度試してみてください。", ctx)
			return err
		}
		s.SendLiveChatMessage(s.ProcessedUserDisplayName + "さんが退室しました！" +
			"（作業時間" + strconv.Itoa(workedTimeSec / 60) + "分）", ctx)
		//s.SendLiveChatMessage(s.ProcessedUserDisplayName + " has finished working! " +
		//	"(" + strconv.Itoa(workedTimeSec / 60) + " minutes)", ctx)
	}
	// ログ記録
	err = s.FirestoreController.AddUserHistory(s.ProcessedUserId, ExitAction, seat, ctx)
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed to add an user history", err)
	}
	// 退室時刻を記録
	err = s.FirestoreController.SetLastExitedDate(s.ProcessedUserId, ctx)
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed to update last-exited-date", err)
		return err
	}
	// 累計学習時間を更新
	err = s.UpdateTotalWorkTime(workedTimeSec, ctx)
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed to update total study time", err)
		return err
	}
	return nil
}

func (s *System) CurrentSeatId(ctx context.Context) (int, customerror.CustomError) {
	// ますは Default room にいるかどうか
	defaultRoomData, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return 0, customerror.Unknown.Wrap(err)
	}
	for _, seat := range defaultRoomData.Seats {
		if seat.UserId == s.ProcessedUserId {
			return seat.SeatId, customerror.NewNilCustomError()
		}
	}
	// default room にいなければ、no-seat-room　にいるかどうか
	noSeatRoomData, err := s.FirestoreController.RetrieveNoSeatRoom(ctx)
	if err != nil {
		return 0, customerror.Unknown.Wrap(err)
	}
	for _, seat := range noSeatRoomData.Seats {
		if seat.UserId == s.ProcessedUserId {
			return 0, customerror.NewNilCustomError()
		}
	}
	// default-roomにもno-seat-roomにもいない
	return -1, customerror.UserNotInAnyRoom.New("the user is not in any room.")
}

func (s *System) IsSeatExist(seatId int, ctx context.Context) (bool, error) {
	// room-layoutを読み込む
	roomLayout, err := s.FirestoreController.RetrieveDefaultRoomLayout(ctx)
	if err != nil {
		return false, err
	}
	for _, seat := range roomLayout.Seats {
		if seat.Id == seatId {
			return true, nil
		}
	}
	return false, nil
}

func (s *System) UpdateTotalWorkTime(workedTimeSec int, ctx context.Context) error {
	userData, err := s.FirestoreController.RetrieveUser(s.ProcessedUserId, ctx)
	if err != nil {
		return err
	}
	previousTotalSec := userData.TotalStudySec
	previousDailyTotalSec := userData.DailyTotalStudySec
	newTotalSec := previousTotalSec + workedTimeSec
	newDailyTotalSec := previousDailyTotalSec + workedTimeSec
	err = s.FirestoreController.UpdateTotalTime(s.ProcessedUserId, newTotalSec, newDailyTotalSec, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *System) IfUserRegistered(ctx context.Context) (bool, error) {
	_, err := s.FirestoreController.RetrieveUser(s.ProcessedUserId, ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (s *System) TotalStudyTimeStrings(ctx context.Context) (string, string, error) {
	userData, err := s.FirestoreController.RetrieveUser(s.ProcessedUserId, ctx)
	if err != nil {
		return "", "", err
	}
	// 累計
	var totalStr string
	totalDuration := time.Duration(userData.TotalStudySec) * time.Second
	if totalDuration < time.Hour {
		totalStr = strconv.Itoa(int(totalDuration.Minutes())) + "分"
		//totalStr = strconv.Itoa(int(totalDuration.Minutes())) + " minutes"
	} else {
		totalStr = strconv.Itoa(int(totalDuration.Hours())) + "時間" +
			strconv.Itoa(int(totalDuration.Minutes()) % 60) + "分"
		//totalStr = strconv.Itoa(int(totalDuration.Hours())) + " hours " +
		//	strconv.Itoa(int(totalDuration.Minutes()) % 60) + " minutes"
	}
	// 当日の累計
	var dailyTotalStr string
	dailyTotalDuration := time.Duration(userData.DailyTotalStudySec) * time.Second
	if dailyTotalDuration < time.Hour {
		dailyTotalStr = strconv.Itoa(int(dailyTotalDuration.Minutes())) + "分"
		//dailyTotalStr = strconv.Itoa(int(dailyTotalDuration.Minutes())) + " minutes"
	} else {
		dailyTotalStr = strconv.Itoa(int(dailyTotalDuration.Hours())) + "時間" +
			strconv.Itoa(int(dailyTotalDuration.Minutes())) + "分"
		//dailyTotalStr = strconv.Itoa(int(dailyTotalDuration.Hours())) + " hours " +
		//	strconv.Itoa(int(dailyTotalDuration.Minutes())) + " minutes"
	}
	return totalStr, dailyTotalStr, nil
}

func (s *System) ExitAllUserDefaultRoom(ctx context.Context) error {
	defaultRoom, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return err
	}
	for _, seat := range defaultRoom.Seats {
		s.ProcessedUserId = seat.UserId
		s.ProcessedUserDisplayName = seat.UserDisplayName
		err := s.ExitRoom(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *System) SendLiveChatMessage(message string, ctx context.Context) {
	err := s.LiveChatBot.PostMessage(message, ctx)
	if err != nil {
		_ = s.LineBot.SendMessageWithError("failed to send live chat message", err)
	}
	return
}

func (s *System) OrganizeDatabase(ctx context.Context) error {
	// untilを過ぎているdefaultルーム内のユーザーを退室させる
	defaultRoom, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return err
	}
	for _, seat := range defaultRoom.Seats {
		if seat.Until.Before(time.Now()) {
			s.ProcessedUserId = seat.UserId
			s.ProcessedUserDisplayName = seat.UserDisplayName
			err := s.ExitRoom(ctx)
			if err != nil {
				return err
			}
		}
	}
	
	// no-seat-roomも同様。
	noSeatRoom, err := s.FirestoreController.RetrieveDefaultRoom(ctx)
	if err != nil {
		return err
	}
	for _, seat := range noSeatRoom.Seats {
		if seat.Until.Before(time.Now()) {
			s.ProcessedUserId = seat.UserId
			s.ProcessedUserDisplayName = seat.UserDisplayName
			err := s.ExitRoom(ctx)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (s *System) ResetDailyTotalStudyTime(ctx context.Context) error {
	log.Println("ResetDailyTotalStudyTime()")
	constantsConfig, err := s.FirestoreController.RetrieveSystemConstantsConfig(ctx)
	if err != nil {
		return err
	}
	previousDate := constantsConfig.LastResetDailyTotalStudySec.Local()
	now := time.Now()
	isDifferentDay := now.Year() != previousDate.Year() || now.Month() != previousDate.Month() || now.Day() != previousDate.Day()
	if isDifferentDay && now.After(previousDate) {
		userRefs, err := s.FirestoreController.RetrieveAllUserDocRefs(ctx)
		if err != nil {
			return err
		}
		for _, userRef := range userRefs {
			err := s.FirestoreController.ResetDailyTotalStudyTime(userRef, ctx)
			if err != nil {
				return err
			}
		}
		log.Println("successfully reset all user's daily total study time.")
		err = s.FirestoreController.SetLastResetDailyTotalStudyTime(now, ctx)
		if err != nil {
			return err
		}
	} else {
		log.Println("all user's daily total study times are already reset today.")
	}
	return nil
}
