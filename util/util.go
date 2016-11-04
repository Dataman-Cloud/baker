package util

func newClient(c *cli.Context) (client.Client, error) {
	var token = c.GlobalString("token")
	var server = c.GlobalString("server")

	// if no server url is provided we can default
	// to the hosted Drone service.
	if len(server) == 0 {
		return nil, fmt.Errorf("Error: you must provide the Drone server address.")
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("Error: you must provide your Drone access token.")
	}

	// attempt to find system CA certs
	certs := syscerts.SystemRootsPool()
	tlsConfig := &tls.Config{RootCAs: certs}

	// create the drone client with TLS options
	return client.NewClientTokenTLS(server, token, tlsConfig), nil
}

// helper function to convert a string slice to a map.
func sliceToMap(s []string) map[string]bool {
	v := map[string]bool{}
	for _, ss := range s {
		v[ss] = true
	}
	return v
}
