package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/httpx"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
	"strconv"
	"strings"
)

type Base struct {
	Logger    *logan.Entry
	Request   *http.Request
	Writer    http.ResponseWriter
	HistoryQ  history.QInterface
	CoreQ     core.QInterface
	PageQuery db2.PageQueryV2
	Filters   map[string]string
}

func (b *Base) Prepare(r *http.Request, w http.ResponseWriter) {
	b.HistoryQ = ctx.HistoryQ(r)
	b.CoreQ = ctx.CoreQ(r)
	b.Logger = ctx.Log(r)
	b.Request = r
	b.Writer = w
}

func (b *Base) GetUint64(name string) (uint64, error) {
	str := b.Request.URL.Query().Get(name)

	if len(str) == 0 {
		return 0, nil
	}

	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Invalid "+name+" query param: "+str)
	}

	return res, nil
}

func (b *Base) GetPageQuery() error {
	limit, err := b.GetUint64("limit")
	if err != nil {
		return err
	}
	page, err := b.GetUint64("page")
	if err != nil {
		return err
	}

	pageQuery, err := db2.NewPageQueryV2(page, limit)
	if err != nil {
		return errors.Wrap(err, "Failed to init page query")
	}

	b.PageQuery = pageQuery
	return nil
}

func (b *Base) IsRequested(inclusion string) bool {
	param := chi.URLParam(b.Request, "include")
	return strings.Contains(param, inclusion)
}

func (b *Base) GetLinksObject() resource.LinksObject {
	path := b.Request.URL.Path
	// TODO: hal.linkBuilder is still useful, move it somewhere
	lb := hal.LinkBuilder{ Base: httpx.BaseURL(b.Request.Context()) }
	format := path + "?page=%d&limit=%d"

	return resource.LinksObject{
		Next: lb.Linkf(b.Filters, format, b.PageQuery.Page + 1, b.PageQuery.Limit).Href,
	}
}

func (b *Base) Render(resource interface{}) {
	js, err := json.MarshalIndent(resource, "", "	")
	if err != nil {
		b.RenderErr(err)
		return
	}

	_, err = b.Writer.Write(js)
	if err != nil {
		b.RenderErr(err)
		return
	}
}

func (b *Base) RenderErr(err error) {
	switch e := err.(type) {
	case *jsonapi.ErrorObject:
		ape.RenderErr(b.Writer, e)
		break
	default:
		b.Logger.Error(err)
		ape.RenderErr(b.Writer, problems.InternalError())
	}
}
