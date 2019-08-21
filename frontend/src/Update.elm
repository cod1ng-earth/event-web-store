module Update exposing (update)

import CartPage.Update
import CatalogPage.Model
import CatalogPage.Update
import ProductDetailPage.Model
import ProductDetailPage.Update
import Model exposing (Content(..), Model)
import Message exposing (Msg(..))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model.content ) of
        ( ProductDetailPageMsg subMsg, ProductDetailPage mdl ) ->
            let
                ( updatedModel, cmd ) =
                    ProductDetailPage.Update.update subMsg mdl
            in
                ( { model | content = ProductDetailPage updatedModel }, cmd )

        ( CatalogPageMsg subMsg, CatalogPage mdl ) ->
            let
                ( updatedModel, cmd ) =
                    CatalogPage.Update.update subMsg mdl
            in
                ( { model | content = CatalogPage updatedModel }, cmd )

        ( CartPageMsg subMsg, _ ) ->
            let
                ( updatedModel, cmd ) =
                    CartPage.Update.update subMsg model.cart
            in
                ( { model | cart = updatedModel }, cmd )

        ( ShowCartPageMsg, _ ) ->
            ( { model | content = ShowCartPage }, Cmd.none )

        ( ShowCatalogPage, _ ) ->
            let
                ( catalogModel, catalogCmd ) =
                    CatalogPage.Model.init
            in
                ( { model | content = CatalogPage catalogModel }, catalogCmd )

        ( ShowProductDetailPage id, _ ) ->
            let
                ( updatedModel, cmd ) =
                    ProductDetailPage.Model.init id
            in
                ( { model | content = ProductDetailPage updatedModel }, cmd )

        ( _, _ ) ->
            -- this means a msg was send for a currently not active content
            ( model, Cmd.none )
