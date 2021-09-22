package bff

import (
	"context"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"

	mGameEngine "github.com/prakharmaurya/m-game-engine/api"
	mHighScore "github.com/prakharmaurya/m-highscore/api"
	"google.golang.org/grpc"
)

type gameResource struct {
	gameClient       mHighScore.GameClient
	gameEngineClient mGameEngine.GameEngineClient
}

func NewGameResource(gameClient mHighScore.GameClient, gameEngineClient mGameEngine.GameEngineClient) *gameResource {
	return &gameResource{
		gameClient:       gameClient,
		gameEngineClient: gameEngineClient,
	}
}

func NewGrpcGameServiceClient(serverAddr string) (mHighScore.GameClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial : %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-highscore connection is nil in m-bff")
	}

	client := mHighScore.NewGameClient(conn)

	return client, nil
}

func NewGrpcGameEngineServiceClient(serverAddr string) (mGameEngine.GameEngineClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial : %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-game-engine connection is nil in m-bff")
	}

	client := mGameEngine.NewGameEngineClient(conn)

	return client, nil
}

func (gr *gameResource) SetHighScore(c *gin.Context) {
	highScoreString := c.Param("hs")
	highScoreFloat64, err := strconv.ParseFloat(highScoreString, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert highscore float")
	}

	gr.gameClient.SetHighScore(context.Background(), &mHighScore.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})
}

func (gr *gameResource) GetHighScore(c *gin.Context) {
	highScoreResponse, err := gr.gameClient.GetHighScore(context.Background(), &mHighScore.GetHighScoreRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting high score")
	}
	hsString := strconv.FormatFloat(highScoreResponse.HighScore, 'e', -1, 64)
	c.JSONP(200, gin.H{
		"hs": hsString,
	})
}

func (gr *gameResource) GetSize(c *gin.Context) {
	sizeResponse, err := gr.gameEngineClient.GetSize(context.Background(), &mGameEngine.GetSizeRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting size")
	}
	c.JSONP(200, gin.H{
		"size": sizeResponse.GetSize(),
	})
}

func (gr *gameResource) SetScore(c *gin.Context) {
	scoreString := c.Param("score")
	scoreFloat64, err := strconv.ParseFloat(scoreString, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert score to float")
	}

	_, err = gr.gameEngineClient.SetScore(context.Background(), &mGameEngine.SetScoreRequest{
		Score: scoreFloat64,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while setting score in m-game-engine")
	}
}
