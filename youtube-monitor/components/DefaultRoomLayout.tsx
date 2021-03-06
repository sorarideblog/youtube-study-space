import React, { useState } from "react";
import styles from "./DefaultRoomLayout.module.sass";
import { RoomLayout } from "../types/room-layout";
import { DefaultRoomState, seat } from "../types/room-state";

class DefaultRoomLayout extends React.Component<
  { layout: RoomLayout; roomState: DefaultRoomState },
  any
> {
  seatWithSeatId(seatId: number, seats: seat[]) {
    let targetSeat: seat = seats[0];
    seats.forEach((seat) => {
      if (seat.seat_id === seatId) {
        targetSeat = seat;
      }
    });
    return targetSeat;
  }

  render() {
    if (this.props.layout && this.props.roomState) {
      const usedSeatIds = this.props.roomState.seats.map(
        (seat) => seat.seat_id
      );

      const emptySeatColor = "#fce7d2";
      const filledSeatColor = "#ee989f";

      const roomLayout = this.props.layout;
      const roomShape = {
        widthPx:
          (900 * roomLayout.room_shape.width) / roomLayout.room_shape.height,
        heightPx: 900,
      };

      const seatFontSizePx = roomShape.widthPx * roomLayout.font_size_ratio;

      const seatShape = {
        width:
          (100 * roomLayout.seat_shape.width) / roomLayout.room_shape.width,
        height:
          (100 * roomLayout.seat_shape.height) / roomLayout.room_shape.height,
      };

      const seatPositions = roomLayout.seats.map((seat) => ({
        x: (100 * seat.x) / roomLayout.room_shape.width,
        y: (100 * seat.y) / roomLayout.room_shape.height,
      }));

      const partitionShapes = roomLayout.partitions.map((partition) => {
        const partitionShapes = roomLayout.partition_shapes;
        const shapeType = partition.shape_type;
        let widthPercent;
        let heightPercent;
        for (let i = 0; i < partitionShapes.length; i++) {
          if (partitionShapes[i].name === shapeType) {
            widthPercent =
              (100 * partitionShapes[i].width) / roomLayout.room_shape.width;
            heightPercent =
              (100 * partitionShapes[i].height) / roomLayout.room_shape.height;
          }
        }
        return {
          widthPercent,
          heightPercent,
        };
      });

      const partitionPositions = roomLayout.partitions.map((partition) => ({
        x: (100 * partition.x) / roomLayout.room_shape.width,
        y: (100 * partition.y) / roomLayout.room_shape.height,
      }));

      const seatList = roomLayout.seats.map((seat, index) => {
        const isUsed = usedSeatIds.includes(seat.id);
        const displayName = usedSeatIds.includes(seat.id)
          ? this.seatWithSeatId(seat.id, this.props.roomState.seats)
              .user_display_name
          : "";
        return (
          <div
            key={seat.id}
            className={styles.seat}
            style={{
              backgroundColor: isUsed ? filledSeatColor : emptySeatColor,
              left: seatPositions[index].x + "%",
              top: seatPositions[index].y + "%",
              width: seatShape.width + "%",
              height: seatShape.height + "%",
              fontSize: seatFontSizePx + "px",
            }}
          >
            <p className={styles.seatNum} style={{ fontWeight: "bold" }}>
              {seat.id}
            </p>
            <p
              className={styles.userDisplayName}
              style={{ overflow: "hidden" }}
            >
              {displayName.substr(0, 4)}
            </p>
          </div>
        );
      });

      const partitionList = roomLayout.partitions.map((partition, index) => (
        <div
          key={partition.id}
          className={styles.partition}
          style={{
            left: partitionPositions[index].x + "%",
            top: partitionPositions[index].y + "%",
            width: partitionShapes[index].widthPercent + "%",
            height: partitionShapes[index].heightPercent + "%",
          }}
        />
      ));

      return (
        <div>
          <div
            id={styles.roomLayout}
            style={{
              width: roomShape.widthPx + "px",
              height: roomShape.heightPx + "px",
            }}
          >
            {seatList}

            {partitionList}
          </div>
        </div>
      );
    } else {
      return <div>Loading</div>;
    }
  }
}

export default DefaultRoomLayout;
