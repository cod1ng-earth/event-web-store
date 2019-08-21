module CatalogPage.Model exposing (Model, init)

import Task
import Time
import Catalog
import Message exposing (Msg)
import CatalogPage.Message


type alias Model =
    { products : Maybe (List Catalog.Product)
    , currentPage : Int
    , totalPages : Int
    , sorting : String
    , prefix : String
    , filtering : String
    , error : Maybe String
    }


init : ( Model, Cmd Msg )
init =
    let
        model =
            { products = Nothing
            , currentPage = 0
            , totalPages = 0
            , sorting = "name"
            , prefix = ""
            , filtering = ""
            , error = Nothing
            }
    in
        ( model, Task.perform (\_ -> Message.CatalogPageMsg CatalogPage.Message.LoadProducts) Time.here )
