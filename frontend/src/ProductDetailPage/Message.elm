module ProductDetailPage.Message exposing (..)

import Http
import Catalog


type SubMsg
    = LoadProduct String
    | GotProduct (Result Http.Error Catalog.Product)



--    | AddToCart String
