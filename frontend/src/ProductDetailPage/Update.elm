module ProductDetailPage.Update exposing (..)

import Http
import Protobuf.Decode as Decode
import Catalog
import Message exposing (..)
import ProductDetailPage.Model exposing (Model)
import ProductDetailPage.Message exposing (..)


update : SubMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LoadProduct uuid ->
            ( model, fetchProduct uuid )

        GotProduct result ->
            case result of
                Ok p ->
                    ( { model | product = Just p, error = Nothing }, Cmd.none )

                Err e ->
                    ( { model | product = Nothing, error = Just (toString e) }, Cmd.none )


fetchProduct : String -> Cmd Msg
fetchProduct id =
    Http.get
        { url = "http://localhost:8080/product?uuid=" ++ id
        , expect = Decode.expectBytes GotProduct Catalog.productDecoder
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
