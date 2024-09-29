package crawler

import (
	"github.com/spf13/viper"
	"news_scrapper/pkg/config"
	"news_scrapper/pkg/crawler/news/queryl"
	"news_scrapper/pkg/style"
)

type CrawlerHelper struct {
	Settings config.Settings
}

func CreateCrawlerHelper(settings config.Settings) CrawlerHelper {
	return CrawlerHelper{Settings: settings}
}
func (ch CrawlerHelper) FetchNews(newsPortalId string) {
	cfg := ch.SearchNewsPortalId(newsPortalId)
	if cfg != nil {
		if cfg.Type == config.Queryly {
			ch.CreateFakingTheFunk(*cfg, func(handler queryl.Handler) {

				if viper.GetString(config.ActionQueryValue) != "" {
					handler.Scrape(queryl.Query,
						viper.GetString(config.ActionQueryValue),
						viper.GetInt(config.ActionLimitSortByDate),
						viper.GetInt(config.ActionLimitBatchSize))
				} else {
					handler.Scrape(queryl.Section,
						viper.GetString(config.ActionSectionFlag),
						viper.GetInt(config.ActionLimitSortByDate),
						viper.GetInt(config.ActionLimitBatchSize))
				}
			})
		} else {

		}
	}
}

func (ch CrawlerHelper) DisplayAvailableOptions(newsPortalId string) {
	cfg := ch.SearchNewsPortalId(newsPortalId)
	if cfg != nil {
		if cfg.Type == config.Queryly {
			ch.CreateFakingTheFunk(*cfg, func(handler queryl.Handler) {
				meta, err := handler.FetchMeta()
				if err != nil {
					return
				}
				style.OnActionF("[%v] Displaying available options now...", newsPortalId)

				style.OnActionF("[%v] ========== Options ", newsPortalId)
				for i := range meta.Section {
					style.OnActionF("[%v] - %v", i, meta.Section[i])
				}
				style.OnActionF("[%v] ========== Dates ", newsPortalId)
				for i := range meta.PubDate {
					style.OnActionF("[%v] - %v", i, meta.PubDate[i])
				}
			})
		} else {

		}
	}
}

func (ch CrawlerHelper) SearchNewsPortalId(newsPortalId string) *config.PortalOption {
	for i := range ch.Settings.Roboto.PortalRegistry {
		portalRegistry := ch.Settings.Roboto.PortalRegistry[i]
		if portalRegistry.Enabled && newsPortalId == portalRegistry.Context {
			style.SuccessfulActionF("[%v] Config found", newsPortalId)
			return &portalRegistry
		}
	}
	style.FailedActionF("Not found id in the registry: %v", newsPortalId)
	return nil
}

func (ch CrawlerHelper) CreateFakingTheFunk(option config.PortalOption, cb func(handler queryl.Handler)) {
	var batch []queryl.FakingTheFunk
	batch = append(batch, queryl.FakingTheFunk{
		option.Enabled, option.Context, option.Host, option.ApiKey, option.Selector, option.OverrideExistingNews,
	})
	queryl.TriggerBatchScrapping(batch, cb)
}

func (ch CrawlerHelper) DisplayAllIds() []string {
	var ids []string
	for i := range ch.Settings.Roboto.PortalRegistry {
		portalRegistry := ch.Settings.Roboto.PortalRegistry[i]
		if portalRegistry.Enabled {
			style.LogF(style.Mounting, "[%v] %s", i, portalRegistry.Context)
			ids = append(ids, portalRegistry.Context)
		}
	}
	return ids
}
