package monitor

import (
	"container-monitor/internal/consts"
	"container-monitor/internal/model"
	"container-monitor/internal/service"
	"context"
	"fmt"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"io"
	"log/slog"
)

var (
	warningWebhook string // message send webhook
	receiverType   string // email sms etc
	receiver       string // warning message receiver
	subject        string // subject
)

type sContainerMonitor struct {
}

func init() {
	service.RegisterContainerMonitor(newContainerMonitor())
}

func newContainerMonitor() *sContainerMonitor {
	return &sContainerMonitor{}
}

func (s *sContainerMonitor) Start() {
	ctx := gctx.New()

	// read ENV first then config
	warningWebhook = genv.Get("WEBHOOK").String()
	receiverType = genv.Get("RECEIVER_TYPE").String()
	receiver = genv.Get("RECEIVER").String()
	subject = genv.Get("SUBJECT").String()
	slog.InfoContext(ctx, "env", "WEBHOOK",
		warningWebhook, "RECEIVER_TYPE", receiverType,
		"RECEIVER", receiver)

	if warningWebhook == "" {
		warningWebhook = g.Cfg().MustGet(ctx, "warn.webhook").String()
	}
	if receiverType == "" {
		receiverType = g.Cfg().MustGet(ctx, "warn.receiverType").String()
	}
	if subject == "" {
		subject = g.Cfg().MustGet(ctx, "warn.subject").String()
	}
	if receiver == "" {
		receiver = g.Cfg().MustGet(ctx, "warn.receiver").String()
	}

	jobName := "ContainerMonitor"

	//Second Minute Hour Day Month Week
	gcron.AddSingleton(ctx, "0 */1 * * * *", func(ctx context.Context) {
		c := gctx.WithCtx(ctx)
		start := gtime.Now()
		err := s.Monitor(c, start)
		if err != nil {
			slog.ErrorContext(c, jobName, "Monitor err", err)
		}
	}, jobName)
}

func (s *sContainerMonitor) Monitor(ctx context.Context, t *gtime.Time) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.ErrorContext(ctx, "Monitor", "NewClientWithOpts err", err)
		return err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{All: true})
	if err != nil {
		slog.ErrorContext(ctx, "Monitor", "ContainerList err", err)
		return err
	}

	warningMsg := ""
	for _, container := range containers {
		//slog.InfoContext(ctx, "Monitor container", "name", container.Names[0], "status", container.Status, "state", container.State)
		if container.State != consts.ContainerStateRunning {
			warningMsg += fmt.Sprintf("container name:%s state:%s\n", container.Names[0], container.State)
		}
	}
	if len(warningMsg) > 0 {
		// send warning msg
		if warningWebhook != "" && receiver != "" && receiverType != "" {
			msg := model.MsgReq{
				Subject: subject,
				MsgType: consts.MsgTypeEmail,
				Content: model.MsgContent{
					Text: warningMsg,
				},
				To: receiver,
			}
			resp, err := g.Client().ContentJson().Post(ctx, warningWebhook, msg)
			if err != nil {
				slog.ErrorContext(ctx, "Monitor", "Post err", err, "msg", msg)
				return err
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.ErrorContext(ctx, "Monitor", "ReadAll err", err)
				return err
			}
			slog.InfoContext(ctx, "Monitor", "Webhook resp data", gconv.String(data))
		} else {
			slog.WarnContext(ctx, "Monitor", "Warning msg", warningMsg)
		}
	}
	fmt.Println("--------------------------------------------------------------------------------------------------")
	return nil
}
