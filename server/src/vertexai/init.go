package vertexai

import (
	"project/util"
	"project/zj"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"github.com/zhengkai/life-go"
	"google.golang.org/api/option"
)

func init() {

	var err error
	theClient, err = aiplatform.NewPredictionClient(
		life.CTX,
		option.WithEndpoint(`us-central1-aiplatform.googleapis.com:443`),
		option.WithCredentialsFile(util.Static(`aigc-llm-730bb179e13c.json`)),
	)
	if err != nil {
		zj.W(err)
	}
}
