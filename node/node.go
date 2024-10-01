package node

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kellabyte/holocron/store"
	"github.com/rs/zerolog"
)

const (
	RoleSolo = iota
	RoleLeader
	RoleFollower
)

var roleNames = make(map[int]string, 0)

type Node struct {
	NodeId uuid.UUID
	epoch  int
	roleId int
	store  *store.Store
	logger zerolog.Logger
}

func New(nodeId uuid.UUID, store *store.Store, logger zerolog.Logger) *Node {
	// TODO: This may not be the right place to initialize this and not thread safe.
	if len(roleNames) == 0 {
		roleNames[RoleSolo] = "solo"
		roleNames[RoleLeader] = "leader"
		roleNames[RoleFollower] = "follower"
	}

	node := &Node{
		NodeId: nodeId,
		store:  store,
		roleId: RoleSolo,
	}
	node.logger = logger.With().Int("epoch", node.epoch).Str("role", node.RoleName()).Logger()
	return node
}

func (node *Node) NodeIdShort() string {
	segments := strings.Split(node.NodeId.String(), "-")
	return segments[4]
}

func NodeIdShort(nodeId uuid.UUID) string {
	segments := strings.Split(nodeId.String(), "-")
	return segments[4]
}

func (node *Node) RoleName() string {
	role, found := roleNames[node.roleId]
	if found {
		return role
	} else {
		return "unknown"
	}
}

func (node *Node) SetRole(roleId int) {
	node.roleId = roleId
}

func (node *Node) Start(ctx context.Context) {
	node.logger.Info().Msg("Starting node")
	go node.startLeadershipPolling(ctx)
}

func (node *Node) GetLeadership() error {
	epoch, err := node.store.GetLatestEpoch()
	if err != nil {
		return err
	}
	epoch++
	node.epoch = epoch

	err = node.store.PutEpoch(epoch, node.NodeId)
	if err != nil {
		node.SetRole(RoleFollower)

		node.logger.Info().
			Int("epoch", node.epoch).
			Str("role", node.RoleName()).
			Msg("Failed to acquire cluster leadership, becoming a follower.")

		return err
	}

	node.SetRole(RoleLeader)
	node.logger.Info().
		Int("epoch", node.epoch).
		Str("role", node.RoleName()).
		Msg("Acquired cluster leadership.")

	return nil
}

func (node *Node) startLeadershipPolling(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		node.logger.Info().
			Int("epoch", node.epoch).
			Str("role", node.RoleName()).
			Msg("Trying to acquire cluster leadership.")

		node.GetLeadership()
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
