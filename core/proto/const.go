package proto

type Capability struct{}

var (
	LONG_PASSWORD                  int = 0x00000001
	FOUND_ROWS                     int = 0x00000002
	LONG_FLAG                      int = 0x00000004
	CONNECT_WITH_DB                int = 0x00000008
	NO_SCHEMA                      int = 0x00000010
	PROTOCOL_41                    int = 0x00000200
	TRANSACTIONS                   int = 0x00002000
	SECURE_CONNECTION              int = 0x00008000
	PLUGIN_AUTH                    int = 0x00080000
	CONNECT_ATTRS                  int = 0x00100000
	PLUGIN_AUTH_LENENC_CLIENT_DATA int = 0x00200000
	SESSION_TRACK                  int = 0x00800000
	DEPRECATE_EOF                  int = 0x01000000
)
