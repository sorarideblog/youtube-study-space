package core

import (
	"app.modules/core/i18n"
	"app.modules/core/myfirestore"
	"app.modules/core/utils"
	"github.com/pkg/errors"
)

func (s *System) ValidateCommand(command utils.CommandDetails) string {
	switch command.CommandType {
	case utils.In:
		return s.ValidateIn(command)
	case utils.Out:
		return ""
	case utils.Info:
		return s.ValidateInfo(command)
	case utils.My:
		return s.ValidateMy(command)
	case utils.Change:
		// seatStateに依存するためChange()の中で行う。
		return ""
	case utils.Seat:
		return s.ValidateSeat(command)
	case utils.Report:
		return s.ValidateReport(command)
	case utils.Kick:
		return s.ValidateKick(command)
	case utils.Check:
		return s.ValidateCheck(command)
	case utils.Block:
		return s.ValidateBlock(command)
	case utils.More:
		return s.ValidateMore(command)
	case utils.Break:
		return s.ValidateBreak(command)
	case utils.Resume:
		return s.ValidateResume(command)
	case utils.Rank:
		return ""
	default:
		return ""
	}
}

func (s *System) ValidateIn(command utils.CommandDetails) string {

	// 作業時間の値
	inputWorkMin := command.InOption.MinutesAndWorkName.DurationMin
	if inputWorkMin != 0 {
		expect := s.Configs.Constants.MinWorkTimeMin <= inputWorkMin && inputWorkMin <= s.Configs.Constants.MaxWorkTimeMin
		if !expect {
			return i18n.T("validate:invalid-work-time-range", s.Configs.Constants.MinWorkTimeMin, s.Configs.Constants.MaxWorkTimeMin)
		}
	}
	// 席番号
	if command.InOption.IsSeatIdSet {
		if command.InOption.SeatId < 0 {
			return i18n.T("validate:negative-seat-id")
		}
	}

	// 作業名は特に制限はない
	// pass

	return ""
}

func (s *System) ValidateInfo(_ utils.CommandDetails) string {
	// pass

	return ""
}

func (s *System) ValidateMy(command utils.CommandDetails) string {
	var isRankVisibleSet, isDefaultStudyMinSet, isFavoriteColorSet bool

	for _, option := range command.MyOptions {
		switch option.Type {
		case utils.RankVisible:
			if isRankVisibleSet {
				return "more than 2 RankVisible options."
			}
			isRankVisibleSet = true
		case utils.DefaultStudyMin:
			if isDefaultStudyMinSet {
				return "more than 2 DefaultStudyMin options."
			}
			inputDefaultStudyMin := option.IntValue
			if inputDefaultStudyMin != 0 {
				expect := s.Configs.Constants.MinWorkTimeMin <= inputDefaultStudyMin && inputDefaultStudyMin <= s.Configs.Constants.MaxWorkTimeMin
				if !expect {
					return i18n.T("validate:invalid-work-time-range", s.Configs.Constants.MinWorkTimeMin, s.Configs.Constants.MaxWorkTimeMin)
				}
			}
			isDefaultStudyMinSet = true
		case utils.FavoriteColor:
			if isFavoriteColorSet {
				return "more than 2 FavoriteColor options."
			}
			expect := utils.IsIncludedInColorNames(option.StringValue) || option.StringValue == ""
			if !expect {
				return i18n.T("validate:invalid-favorite-color-option", utils.FavoriteColorMyOptionPrefix)
			}
			isFavoriteColorSet = true
		default:
			return "there is an unknown option in command.MyOptions"
		}
	}
	return ""
}

func (s *System) ValidateSeat(_ utils.CommandDetails) string {
	// pass

	return ""
}

func (s *System) ValidateKick(command utils.CommandDetails) string {
	// 指定座席番号
	if command.KickOption.SeatId <= 0 {
		return i18n.T("validate:non-one-or-more-seat-id")
	}

	return ""
}

func (s *System) ValidateCheck(command utils.CommandDetails) string {
	// 指定座席番号
	if command.CheckOption.SeatId <= 0 {
		return i18n.T("validate:non-one-or-more-seat-id")
	}

	return ""
}

func (s *System) ValidateBlock(command utils.CommandDetails) string {
	// 指定座席番号
	if command.BlockOption.SeatId <= 0 {
		return "席番号は1以上にしてください。"
	}

	return ""
}

func (s *System) ValidateReport(command utils.CommandDetails) string {
	// 空欄でないか
	if command.ReportOption.Message == "" {
		return i18n.T("validate:parse:missing-message", utils.ReportCommand)
	}

	return ""
}

func (s *System) ValidateChange(command utils.CommandDetails, seatState myfirestore.SeatState) error {
	// オプションが1つ以上指定されているか
	if command.ChangeOption.NumOptionsSet() == 0 {
		return errors.New(i18n.T("validate:missing-option"))
	}

	switch seatState {
	case myfirestore.WorkState:
		// 作業内容
		// pass

		// 入室時間
		if command.ChangeOption.IsDurationMinSet {
			inputDurationMin := command.ChangeOption.DurationMin
			expect := s.Configs.Constants.MinWorkTimeMin <= inputDurationMin && inputDurationMin <= s.Configs.Constants.MaxWorkTimeMin
			if !expect {
				return errors.New(i18n.T("validate:invalid-work-time-range", s.Configs.Constants.MinWorkTimeMin, s.Configs.Constants.MaxWorkTimeMin))
			}
		}
	case myfirestore.BreakState:
		// 休憩内容
		// pass

		// 休憩時間
		if command.ChangeOption.IsDurationMinSet {
			inputDurationMin := command.ChangeOption.DurationMin
			expect := s.Configs.Constants.MinBreakDurationMin <= inputDurationMin && inputDurationMin <= s.Configs.Constants.MaxBreakDurationMin
			if !expect {
				return errors.New(i18n.T("validate:invalid-break-time-range", s.Configs.Constants.MinBreakDurationMin, s.Configs.Constants.MaxBreakDurationMin))
			}
		}
	}

	return nil
}

func (s *System) ValidateMore(command utils.CommandDetails) string {
	// 時間オプション
	if command.MoreOption.DurationMin <= 0 {
		return i18n.T("validate:non-one-or-more-extended-time")
	}

	return ""
}

func (s *System) ValidateBreak(command utils.CommandDetails) string {
	// 休憩内容
	// pass

	// 休憩時間
	if command.BreakOption.IsDurationMinSet {
		inputDurationMin := command.BreakOption.DurationMin
		expect := s.Configs.Constants.MinBreakDurationMin <= inputDurationMin && inputDurationMin <= s.Configs.Constants.MaxBreakDurationMin
		if !expect {
			return i18n.T("validate:invalid-break-time-range", s.Configs.Constants.MinBreakDurationMin, s.Configs.Constants.MaxBreakDurationMin)
		}
	}

	return ""
}

func (s *System) ValidateResume(_ utils.CommandDetails) string {
	// 作業名
	// pass

	return ""
}
