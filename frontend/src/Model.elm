module Model exposing (..)


import Browser.Events exposing (Visibility(..))












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
