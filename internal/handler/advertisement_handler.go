package handler

import (
	"advertisement/internal/convert"
	"advertisement/model"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

func (h *Handler) PlacementRequest(c *gin.Context) {
	var input model.RequestPlacement
	if err := c.BindJSON(&input); err != nil {
		h.log.Error("Error receiving request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, string("Error"))
		return
	}

	errValid := h.valid.Struct(input)
	if errValid != nil {
		h.log.Error("Error 400", zap.Error(errValid))
		c.JSON(http.StatusBadRequest, string("Error"))
		return
	}

	h.BidRequest(input, c)
	c.JSON(http.StatusOK, string("OK"))
}

func (h *Handler) BidRequest(input model.RequestPlacement, c *gin.Context) {
	var models []model.Imp
	for iter := 0; iter < len(input.Tiles); iter++ {
		models = append(models, convert.ConvertTilesToImp(input.Tiles[iter]))
	}

	bid := &model.RequestToPartner{
		ID:      input.ID,
		Imp:     models,
		Context: input.Context,
	}
	bytesRequest, err := json.Marshal(bid)
	if err != nil {
		h.log.Error("Error converting structure to json", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}
	ctx, cansel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cansel()

	c.Request = c.Request.WithContext(ctx)

	wg.Add(len(h.address))
	resultList := []*model.PartnerRequest{}
	for iter := 0; iter < len(h.address); iter++ {
		go func(iter int) {
			defer wg.Done()
			res, _ := h.sendRequestToPartner(h.address[iter], bytesRequest, c)
			resultList = append(resultList, res)
		}(iter)
	}
	wg.Wait()

	item := selectMaxPrice(resultList, bid)
	c.JSON(http.StatusOK, item)
}

func (h *Handler) sendRequestToPartner(address string, bytesRequest []byte, c *gin.Context) (*model.PartnerRequest, error) {
	url := fmt.Sprintf("http://%s/bid_request", address)
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, url, nil)
	if err != nil {
		h.log.Error("Request return error", zap.Error(err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.log.Error("Error sending a request to Advertising partners or receiving a response", zap.Error(err))
		return nil, err
	}

	var input model.PartnerRequest
	if err := json.NewDecoder(resp.Body).Decode(&input); err != nil {
		h.log.Error(
			"Error converting jsion to response structure from advertising partners",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, string("Error"))
		return nil, err
	}
	return &input, nil
}

func selectMaxPrice(resultList []*model.PartnerRequest, bid *model.RequestToPartner) *model.Answer {
	models := &model.Answer{
		ID: bid.ID,
	}
	cache := make(map[uint]model.ImpPart)
	for iter := 0; iter < len(resultList); iter++ {
		item := resultList[iter]
		if item == nil {
			continue
		}
		for iterator := 0; iterator < len(bid.Imp); iterator++ {
			maxItem := model.ImpPart{}
			max := float32(0.0)
			for iterator2 := 0; iterator2 < len(item.ImpPart); iterator2++ {
				if bid.Imp[iterator].ID == item.ImpPart[iterator2].ID {
					if item.ImpPart[iterator].Price > max {
						max = item.ImpPart[iterator].Price
						maxItem = item.ImpPart[iterator2]
					}
				}
			}
			if max != 0 {
				val, ok := cache[maxItem.ID]
				if ok {
					if maxItem.Price > val.Price {
						cache[maxItem.ID] = maxItem
					}
				} else {
					cache[maxItem.ID] = maxItem
				}
			}
		}
	}
	for _, v := range cache {
		models.ImpAns = append(models.ImpAns, convert.ConvertPartnerToAnswer(v))
	}
	return models
}
