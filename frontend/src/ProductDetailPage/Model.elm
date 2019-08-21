module ProductDetailPage.Model exposing (..)

import Catalog
import Message exposing (Msg)


type Model
    = Loading
    | Error String
    | Product Catalog.Product


init : Model
init = Loading
