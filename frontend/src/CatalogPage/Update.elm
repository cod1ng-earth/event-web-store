module CatalogPage.Update exposing (..)

import Http
import Url exposing (percentEncode)
import Protobuf.Decode as Decode
import Catalog
import CatalogPage.Model exposing (Model)
import Message exposing (..)
import CatalogPage.Message exposing (..)


update : SubMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        PreviousPage ->
            let
                updated =
                    { model | currentPage = model.currentPage - 1 }
            in
                ( updated, fetchProducts updated )

        NextPage ->
            let
                updated =
                    { model | currentPage = model.currentPage + 1 }
            in
                ( updated, fetchProducts updated )

        SortByUuid ->
            let
                updated =
                    { model | sorting = "uuid" }
            in
                ( updated, fetchProducts updated )

        SortByPrice ->
            let
                updated =
                    { model | sorting = "price" }
            in
                ( updated, fetchProducts updated )

        SortByName ->
            let
                updated =
                    { model | sorting = "name" }
            in
                ( updated, fetchProducts updated )

        DisableFilterByPrefix ->
            let
                updated =
                    { model | filtering = "" }
            in
                ( updated, fetchProducts updated )

        EnableFilterByPrefix ->
            let
                updated =
                    { model | filtering = "prefix" }
            in
                ( updated, fetchProducts updated )

        GoToPage page ->
            let
                updated =
                    { model | currentPage = (Maybe.withDefault 0 (String.toInt page) - 1) }
            in
                ( updated, fetchProducts updated )

        SetFilterPrefix prefix ->
            let
                updated =
                    { model | prefix = prefix }
            in
                ( updated, fetchProducts updated )

        LoadProducts ->
            ( model, fetchProducts model )

        GotProducts result ->
            case result of
                Ok pp ->
                    let
                        updated =
                            case pp.request of
                                Nothing ->
                                    model
                                Just catalogRequest ->
                                    if model.currentPage /= catalogRequest.page then
                                        model
                                    else
                                        { model | products = Just pp.products, currentPage = catalogRequest.page, totalPages = pp.totalPages, error = Nothing }
                    in
                        ( updated, Cmd.none )

                Err e ->
                    let
                        updated =
                            { model | products = Nothing, error = Just (toString e) }
                    in
                        ( updated, Cmd.none )


fetchProducts : Model -> Cmd Msg
fetchProducts model =
    let
        prefix =
            if model.filtering == "" then
                ""
            else
                percentEncode model.prefix
    in
        Http.get
            { url = "http://localhost:8080/products?itemsPerPage=100&sort=" ++ model.sorting ++ "&prefix=" ++ prefix ++ "&page=" ++ String.fromInt model.currentPage
            , expect = Decode.expectBytes GotProducts Catalog.catalogResponseDecoder
            }
            |> Cmd.map CatalogPageMsg


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
