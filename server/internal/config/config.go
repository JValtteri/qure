package config

import (
	"log"
	"runtime"

	"github.com/JValtteri/qure/server/internal/utils"
)


const CONFIG_FILE string = "config.json"

type Config struct {
	ORIGIN_URL					string
	SERVER_PORT					string
	ENABLE_TLS					bool
	CERT_FILE					string
	PRIVATE_KEY_FILE			string
	SOURCE_DIR					string
	DB_FILE_NAME				string
	MIN_USERNAME_LENGTH			int
	MIN_PASSWORD_LENGTH			int
	MAX_SESSION_AGE				utils.Epoch
	MAX_PENDIG_RESERVATION_TIME	utils.Epoch
	RESERVATION_OVERTIME		utils.Epoch
	TEMP_CLIENT_AGE				utils.Epoch		// undocumented
	FIRST_PASSWORD_LENGTH		int				// undocumented
	SESSION_KEY_LENGTH			int
	HASH_MEMORY					uint32			// HASH settings should not be changed from defaults
	HASH_ITERATIONS				uint32			// Defaults are OWASP recommendations
	HASH_PARALLELISM			uint8
	EXTRA_STRICT_SESSIONS		bool			// Enables detecting session ID counterfits. Defaults to passive blocking counterfit attempts.
	MAX_THREADS					int				// Maximum amount of Go threads allowed
	RATE_LIMIT_PER_MINUTE		float64			// Maximum allowed requests/minute per client
	RATE_LIMIT_PER_MINUTE_EVENT	float64			// Maximum allowed requests/minute per client ()
	RATE_LIMIT_BURST			float64			// Maximum allowed burst (on top of base limit)
	RATE_LIMIT_RESET_MINUTES	float32			// Interval to clear reset limiters (to purge old clients and reset counters)
	RATE_LIMIT_ALERT			uint64			// Exceeding this number of blocked requests triggers an alert in log with offending IP address and blocked request count at last limit reset
}

var CONFIG = Config{
	ORIGIN_URL:						"localhost",
	SERVER_PORT:					"8080",
	ENABLE_TLS:						false,
	CERT_FILE:						"cert.pem",
	PRIVATE_KEY_FILE:				"privkey.pem",
	SOURCE_DIR:						"../client/dist",
	DB_FILE_NAME:					"./db/db.gob",
	MIN_USERNAME_LENGTH:			4,
	MIN_PASSWORD_LENGTH:			8,
	MAX_SESSION_AGE:				60*60*24*30,	// max age in seconds
	MAX_PENDIG_RESERVATION_TIME:	60*10,			// seconds
	RESERVATION_OVERTIME:			60*60,			// seconds a reservation is kept past reservation start time
	TEMP_CLIENT_AGE: 				60*60*24*30,	// max age in seconds
	SESSION_KEY_LENGTH:				20,				// Length of the session key stored in the session cookie
	FIRST_PASSWORD_LENGTH:			25,				// Automatically generated pasword length for the first admin user
	HASH_MEMORY:					19*1024,		// HASH settings should not be changed from defaults
	HASH_ITERATIONS:				2,				// Defaults are OWASP recommendations
	HASH_PARALLELISM:				1,
	EXTRA_STRICT_SESSIONS:			false,			// Active counterfit detection: High resource use
	MAX_THREADS:					0,				// 0 is automatic, set this manually if you encounter performance issues with container
	RATE_LIMIT_PER_MINUTE:			60,
	RATE_LIMIT_PER_MINUTE_EVENT:	120,			// Event data requests are more common, so a higher limit is justified
	RATE_LIMIT_BURST:				5,
	RATE_LIMIT_RESET_MINUTES:		60,
	RATE_LIMIT_ALERT:				1000,
}

func LoadConfig(configName string) {
	raw_config := readConfig(configName)
	unmarshal(raw_config, &CONFIG)
	log.Printf("Server url/port: %s:%s\n", CONFIG.ORIGIN_URL, CONFIG.SERVER_PORT)
	if CONFIG.ENABLE_TLS {
		log.Println("TLS is Enabled")
	} else {
		log.Println("HTTP-Only mode")
	}
	if CONFIG.SESSION_KEY_LENGTH < 20 {				// Do not allow short session keys
		CONFIG.SESSION_KEY_LENGTH = 20
	}
	var osThreads = runtime.GOMAXPROCS(CONFIG.MAX_THREADS)
	log.Printf("Detected %v available OS Threads\n", osThreads)
	osThreads = runtime.GOMAXPROCS(CONFIG.MAX_THREADS)
	log.Printf("Using %v threads\n", osThreads)
}

func readConfig(fileName string) []byte {
	raw_config := utils.ReadFile(fileName)
	return raw_config
}

func unmarshal(data []byte, config any) {
	err := utils.LoadJSON(data, config)
	if err != nil {
		log.Fatalf("JSON unmarshal error: %v" , err)
	}
}
