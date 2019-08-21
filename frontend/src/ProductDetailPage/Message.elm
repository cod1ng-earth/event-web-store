module ProductDetailPage.Message exposing (..)

import Http
import Catalog


type SubMsg
    = LoadProduct String
    | PassedSlowLoadThreshold
    | ProductFetched (Result Http.Error Catalog.Product)
