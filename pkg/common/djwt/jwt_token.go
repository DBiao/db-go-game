package djwt

type JwtToken struct {
	Token  string
	Expire int64
	Uid    int64
}
