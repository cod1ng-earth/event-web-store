module CatalogPage.Model exposing (..)

import Catalog


type alias Model =
    { products : Maybe (List Catalog.Product)
    , currentPage : Int
    , totalPages : Int
    , sorting : String
    , prefix : String
    , filtering : String
    , error : Maybe String
    }


init : Model
init =
    { products = Nothing
    , currentPage = 0
    , totalPages = 0
    , sorting = "name"
    , prefix = ""
    , filtering = ""
    , error = Nothing
    }
