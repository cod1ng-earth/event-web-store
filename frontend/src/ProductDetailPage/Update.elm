module ProductDetailPage.Update exposing (update)

import Http
import Protobuf.Decode as Decode
import Catalog
import Message exposing (Msg(..))
import ProductDetailPage.Model exposing (Model(..))
import ProductDetailPage.Message exposing (SubMsg(..))


update : SubMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model ) of
        ( PassedSlowLoadThreshold, Loading ) ->
            ( LoadingSlowly, Cmd.none )

        ( PassedSlowLoadThreshold, _ ) ->
            ( model, Cmd.none )

        ( ProductFetched result, _ ) ->
            case result of
                Ok p ->
                    ( Loaded p, Cmd.none )

                Err e ->
                    ( Failed (toString e), Cmd.none )

        ( LoadProduct id, _ ) ->
            ( model, fetchProduct id )


fetchProduct : String -> Cmd Msg
fetchProduct id =
    Http.get
        { url = "http://localhost:8080/product?uuid=" ++ id
        , expect = Decode.expectBytes ProductFetched Catalog.productResponseDecoder
        }
        |> Cmd.map ProductDetailPageMsg


toString : Http.Error -> String
toString error =
    case error of
        Http.BadUrl url ->
            url ++ " is bad"

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network Error"

        Http.BadStatus status ->
            String.fromInt status ++ " Status"

        Http.BadBody msg ->
            msg
