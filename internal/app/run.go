package app

func (a *Application) Run() error {
	//// close database session when server terminates
	//if a.Config.DatabaseConfig.Dialect != "" {
	//	defer a.Database.Session.Close()
	//}
	//
	//// close cache connection when server terminates
	//if a.Config.CacheType != "" {
	//	defer a.Cache.Close()
	//}

	//// add badger garbage collection to scheduler
	//if a.Config.CacheType == "badger" {
	//	cid, err := a.Scheduler.AddFunc("@daily", func() { a.Cache.GetConnection().(*badger.DB).RunValueLogGC(0.7) })
	//	if err != nil {
	//		a.Log.Error().Err(err).Msg("failed to a GC for badger")
	//	} else {
	//		a.Log.Info().Str("cronEntryID", fmt.Sprint(cid)).Msg("added GC for badger to cron jobs")
	//	}
	//}

	//// start mail listener
	//if a.Config.MailerService != "" {
	//	a.Log.Info().Msg("starting mail channels")
	//	go a.Mailer.ListenForMail()
	//}

	// start RPC server
	// TODO: disabled
	//a.Log.Info().Str("port", a.RPCServer.Port).Msg("starting RPC server")
	//go a.RPCServer.Run()

	// start app server
	a.Log.Info().Str("port", a.Config.Port).Msg("starting app server")
	if err := a.Server.Run(); err != nil {
		a.Log.Fatal().Err(err).Msg("failed to start app server")
	}

	return nil
}
