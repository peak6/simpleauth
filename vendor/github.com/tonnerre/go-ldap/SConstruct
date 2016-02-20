opts = Variables( 'options.conf', ARGUMENTS )
opts.Add("DESTDIR", 'Set the root directory to install into ( /path/to/DESTDIR )', "")

env = Environment(ENV = {'GOROOT': '/usr/lib/go'}, TOOLS=['default', 'go'],
		  options = opts)

ldap = env.Go('goldap', ["bind.go", "conn.go", "control.go", "filter.go",
			 "ldap.go", "search.go"])
pack = env.GoPack('ldap', ldap)

env.Install(env['DESTDIR'] + env['GO_PKGROOT'] +
	    '/github.com/mmitton', pack)
env.Alias('install', [env['DESTDIR'] + env['GO_PKGROOT'] +
	  '/github.com/mmitton'])

opts.Save('options.conf', env)
