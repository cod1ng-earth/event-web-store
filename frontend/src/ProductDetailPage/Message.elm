module ProductDetailPage.Message exposing (SubMsg(..))

import Http
import Catalog


type SubMsg
    = LoadProduct String
    | PassedSlowLoadThreshold
    | ProductFetched (Result Http.Error Catalog.ProductResponse)
