package endpoint

const Base string = "https://manifold.markets/api"
const Version string = "v0/"

const DefaultLimit = 1000

const GetBets string = "bets/"
const GetComments string = "comments/"
const GetGroupBySlug string = "group/"
const GetGroupByID string = "group/by-id/"
const GetGroups string = "groups/"
const GetMarketBySlug string = "slug/"
const GetMarketByID string = "market/"
const GetMarkets string = "markets/"
const GetSearchMarkets string = "search-markets/"
const GetMe string = "me/"
const GetUserByUsername string = "user/"
const GetUserByID string = "user/by-id/"
const GetUsers string = "users/"

const PostBet string = "bet/"
const PostCancellation string = "bet/cancel/"
const PostMarket string = "market/"
const PostComment string = "comment/"

const MarketsSuffix = "/markets/"
const PositionsSuffix = "/positions/"
const LiquiditySuffix = "/liquidity/"
const ClosureSuffix = "/close/"
const GroupSuffix = "/group/"
const ResolutionSuffix = "/resolve/"
const SellSuffix = "/sell/"
