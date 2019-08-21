module ProductDetailPage.Model exposing (Model(..), init)

import Process
import Task
import Time
import Catalog
import Message exposing (Msg)
import ProductDetailPage.Message exposing (SubMsg(..))


type Model
    = Loading
    | LoadingSlowly
    | Loaded Catalog.Product
    | Failed String


init : String -> ( Model, Cmd Msg )
init id =
    ( Loading
    , Cmd.batch
        [ Task.perform (\_ -> Message.ProductDetailPageMsg (LoadProduct id)) Time.here
        , Task.perform (\_ -> Message.ProductDetailPageMsg PassedSlowLoadThreshold) (Process.sleep 500)
        ]
    )
