module ProductDetailPage.Model exposing (..)

import Catalog
import Message exposing (Msg)


type alias Model =
    { error : Maybe String
    , product : Maybe Catalog.Product
    }


init : Model
init =
    { product = Nothing
    , error = Nothing
    }
