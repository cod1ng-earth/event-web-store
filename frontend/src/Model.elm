module Model exposing (Model, Content(..), init)

import Browser.Events exposing (Visibility(..))
import CatalogPage.Model
import ProductDetailPage.Model
import CartPage.Model
import Message exposing (Msg)


type alias Model =
    { error : String
    , content : Content
    , cart : CartPage.Model.Model
    }


type Content
    = CatalogPage CatalogPage.Model.Model
    | ProductDetailPage ProductDetailPage.Model.Model
    | ShowCartPage
    | OrderSuccessfulPage
    | OrderFailedPage


init : () -> ( Model, Cmd Msg )
init _ =
    let
        ( cartModel, cartCmd ) =
            CartPage.Model.init

        ( catalogModel, catalogCmd ) =
            CatalogPage.Model.init
    in
        ( { error = ""
          , content = CatalogPage catalogModel
          , cart = cartModel
          }
        , Cmd.batch
            [ cartCmd
            , catalogCmd
            ]
        )
