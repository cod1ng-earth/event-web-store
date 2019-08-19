module Model exposing (..)

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1, h2, h3, ol, ul, li, a, span, header, footer, i, img, input)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (href, class, id, disabled, attribute, src, placeholder)
import List exposing (foldl, append)
import Time
import Task
import Random
import Http
import Delay
import Round exposing (round)
import Catalog
import Checkout
import CatalogPage.Update as CatalogUpdate
import CartPage.Update as CartPageUpdate
import CatalogPage.Model as CatalogModel
import ProductDetailPage.Model as ProductDetailPageModel
import CartPage.Model as CartPageModel
import Message exposing (..)


type alias Model =
    { error : String
    , content : Content
    , cart : CartPageModel.Model
    }


type Content
    = CatalogPage CatalogModel.Model
    | ProductDetailPage ProductDetailPageModel.Model
    | ShowCartPage
    | OrderSuccessfulPage
    | OrderFailedPage


init : () -> ( Model, Cmd Msg )
init _ =
    ( { error = ""
      , content = CatalogPage CatalogModel.init
      , cart = CartPageModel.init
      }
    , Cmd.batch
        [ CatalogUpdate.fetchProducts CatalogModel.init
        , CartPageUpdate.fetchCart
        ]
    )
