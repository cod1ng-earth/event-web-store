module CatalogPage.Message exposing (..)

import Http
import Catalog


type SubMsg
    = GotProducts (Result Http.Error Catalog.CatalogPage)
    | LoadProducts
    | PreviousPage
    | NextPage
    | SortByUuid
    | SortByPrice
    | SortByName
    | DisableFilterByPrefix
    | EnableFilterByPrefix
    | GoToPage String
    | SetFilterPrefix String
