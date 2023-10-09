package main

import (
	"fmt"
	"log"
	"os"

	"koushoku/database"
	. "koushoku/services"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Delete    []int64 `long:"delete" description:"Delete archive(s) by id from the database"`
	DeleteAll bool    `long:"delete-all" description:"Delete all archives from the database"`

	Publish    []int64 `long:"publish" description:"Publish archive(s) by id"`
	PublishAll bool    `long:"publish-all" description:"Publish all archives"`

	Unpublish    []int64 `long:"unpublish" description:"Unpublish archive(s) by id"`
	UnpublishAll bool    `long:"unpublish-all" description:"Unpublish all archives"`

	Moderate bool     `long:"moderate" description:"Moderate all archives (blacklist)"`
	Index    bool     `long:"index" description:"Index archives"`
	Reindex  bool     `long:"reindex" description:"Reindex archives"`
	Add      []string `long:"add" description:"Index archive(s) from path"`

	UpdateSlugs bool `long:"update-slugs" description:"Update slugs for all archives"`
	Purge       bool `long:"purge" description:"Purge symlinks"`
	Remap       bool `long:"remap" description:"Remap symlinks"`

	PurgeThumbnails    bool `long:"purge-thumbnails" description:"Purge thumbnails"`
	GenerateThumbnails bool `long:"generate-thumbnails" description:"Generate thumbnails"`

	Scrape     bool    `long:"scrape" description:"Scrape archives metadata from you-know-where"`
	ScrapeById []int64 `long:"scrape-id" description:"Scrape archive(s) metadata by id from you-know-where"`
	Import     bool    `long:"import" description:"Import metadata from metadata.json"`
	Fpath      string  `long:"fpath" description:"F Path to scrape metadata from"`
	IPath      string  `long:"ipath" description:"I Path to scrape metadata from"`

	Accept []int64 `long:"accept" description:"Accept submission(s) by id"`
	Reject []int64 `long:"reject" description:"Reject submission(s) by id"`
	Note   string  `long:"note" description:"Note for the submission"`

	Subs bool    `long:"subs" description:"List submissions"`
	Sub  int64   `long:"sub" description:"Submission id"`
	Link []int64 `long:"link" description:"Link archive(s) by id to a submission"`

	Archives []int64 `long:"archive"`
	Expunge  bool    `long:"expunge"`
	Redirect int64   `long:"redirect"`
	Source   string  `long:"source"`

	StartPort             int  `long:"start-port"`
	EndPort               int  `long:"end-port"`
	PurgeCaches           bool `long:"purge-caches"`
	PurgeArchivesCache    bool `long:"purge-archives-cache"`
	PurgeTaxonomiesCache  bool `long:"purge-taxonomies-cache"`
	PurgeTemplatesCache   bool `long:"purge-templates-cache"`
	PurgeSubmissionsCache bool `long:"purge-submissions-cache"`
	ReloadTemplates       bool `long:"reload-templates"`
}

func main() {
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		if !flags.WroteHelp(err) {
			log.Fatalln(err)
		}
		return
	}
	database.Init()

	if len(opts.Delete) > 0 {
		log.Println("Deleting archives from the database...")
		for _, id := range opts.Delete {
			if err := DeleteArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.DeleteAll {
		log.Println("Deleting all archives from the database...")
		if err := DeleteArchives(); err != nil {
			log.Fatalln(err)
		}
	}

	if opts.Purge {
		log.Println("Purging symlinks...")
		purgeSymlinks()
	}

	if opts.PurgeThumbnails {
		log.Println("Purging thumbnails...")
		purgeThumbnails()
	}

	if len(opts.Add) > 0 {
		log.Println("Indexing archive...")
		for _, path := range opts.Add {
			indexArchive(path, false)
		}
	}

	if opts.Index {
		log.Println("Indexing archives...")
		indexArchives(false)
	} else if opts.Reindex {
		log.Println("Reindexing archives...")
		indexArchives(true)
	}

	if opts.Remap {
		log.Println("Remapping archives...")
		remapArchives()
	}

	if opts.Scrape {
		log.Println("Scraping metadata...")
		scrapeMetadata()
	}

	if len(opts.ScrapeById) > 0 {
		log.Println("Scraping metadata...")
		for _, id := range opts.ScrapeById {
			scrapeMetadataById(id, opts.Fpath, opts.IPath)
		}
	}

	if opts.Import {
		log.Println("Importing metadata...")
		importMetadata()
	}

	if opts.Moderate {
		log.Println("Moderating archives...")
		moderateArchives()
	}

	if len(opts.Accept) > 0 {
		log.Println("Accepting submissions...")
		for _, id := range opts.Accept {
			if err := AcceptSubmission(id, opts.Note); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if len(opts.Reject) > 0 {
		log.Println("Rejecting submissions...")
		for _, id := range opts.Reject {
			if err := RejectSubmission(id, opts.Note); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.Subs {
		submissions, err := ListSubmissions()
		if err != nil {
			log.Fatalln(err)
		}

		for _, submission := range submissions {
			fmt.Printf("%d: %s\n", submission.ID, submission.Name)
			fmt.Printf("Accepted?: %t\n", submission.Accepted)
			fmt.Println("Rejected?:", submission.Rejected)
			fmt.Println(submission.Content)
		}
	}

	if len(opts.Link) > 0 && opts.Sub > 0 {
		log.Println("Linking submissions...")
		for _, id := range opts.Link {
			if err := LinkSubmission(id, opts.Sub); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if len(opts.Publish) > 0 {
		log.Println("Publishing archives...")
		for _, id := range opts.Publish {
			if _, err := PublishArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.PublishAll {
		log.Println("Publishing all archives...")
		if err := PublishArchives(); err != nil {
			log.Fatalln(err)
		}
	}

	if len(opts.Unpublish) > 0 {
		log.Println("Unpublishing archives...")
		for _, id := range opts.Unpublish {
			if _, err := UnpublishArchive(id); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if opts.UnpublishAll {
		log.Println("Unpublishing all archives...")
		if err := UnpublishArchives(); err != nil {
			log.Fatalln(err)
		}
	}

	if opts.GenerateThumbnails {
		log.Println("Generating thumbnails...")
		generateThumbnails()
	}

	if opts.UpdateSlugs {
		log.Println("Updating slugs...")
		updateSlugs()
	}

	if len(opts.Archives) > 0 && (opts.Redirect > 0 || opts.Expunge || len(opts.Source) > 0) {
		for _, id := range opts.Archives {
			if opts.Expunge {
				log.Println("Expunging archive", id)
				ExpungeArchive(id)
			}

			if opts.Redirect > 0 {
				log.Printf("Redirecting archive %d to %d\n", id, opts.Redirect)
				RedirectArchive(id, opts.Redirect)
			}

			if len(opts.Source) > 0 {
				log.Printf("Setting archive %d source to %s\n", id, opts.Source)
				SetArchiveSource(id, opts.Source)
			}
		}
	}

	if opts.PurgeCaches || opts.PurgeArchivesCache || opts.PurgeTaxonomiesCache ||
		opts.PurgeTemplatesCache || opts.PurgeSubmissionsCache {
		log.Println("Purging caches...")
		purgeCaches(opts.StartPort, opts.EndPort, PurgeCacheOptions{
			Archives:    opts.PurgeArchivesCache || opts.PurgeCaches,
			Taxonomies:  opts.PurgeTaxonomiesCache || opts.PurgeCaches,
			Templates:   opts.PurgeTemplatesCache || opts.PurgeCaches,
			Submissions: opts.PurgeSubmissionsCache || opts.PurgeCaches,
		})
	}

	if opts.ReloadTemplates {
		log.Println("Reloading templates...")
		reloadTemplates(opts.StartPort, opts.EndPort)
	}
}
