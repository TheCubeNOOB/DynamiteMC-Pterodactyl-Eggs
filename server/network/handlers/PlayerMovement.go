package handlers

import (
	"github.com/aimjel/minecraft/packet"
	"github.com/dynamitemc/dynamite/server/player"
)

func PlayerMovement(
	controller Controller,
	state *player.Player,
	p packet.Packet,
) {
	if state.IsDead() {
		return
	}
	x, y, z := state.Position()
	yaw, pitch := state.Rotation()
	switch pk := p.(type) {
	case *packet.PlayerPosition:
		{
			controller.BroadcastMovement(pk.ID(), pk.X, pk.FeetY, pk.Z, yaw, pitch, pk.OnGround, false)
			controller.HandleCenterChunk(x, z, pk.X, pk.Z)
		}
	case *packet.PlayerPositionRotation:
		{
			controller.BroadcastMovement(pk.ID(), pk.X, pk.FeetY, pk.Z, pk.Yaw, pk.Pitch, pk.OnGround, false)
			controller.HandleCenterChunk(x, z, pk.X, pk.Z)
		}
	case *packet.PlayerRotation:
		{
			controller.BroadcastMovement(pk.ID(), x, y, z, pk.Yaw, pk.Pitch, pk.OnGround, false)
		}
	}
}
