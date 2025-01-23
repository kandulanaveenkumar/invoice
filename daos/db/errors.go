package db

import "errors"

var (
	ErrTabelAudit = errors.New("error while auditing the change on table")

	ErrFetchAWBInfo  = errors.New("error while fetching airway bill info")
	ErrUpsertAWBInfo = errors.New("error while upserting airway bill info")
	ErrDeleteAWBInfo = errors.New("error while deleting airway bill info")

	ErrFetchAWBRoutes  = errors.New("error while fetching airway bill routes")
	ErrUpsertAWBRoutes = errors.New("error while upserting airway bill routes")
	ErrDeleteAWBRoutes = errors.New("error while deleting airway bill routes")

	ErrFetchAWBCharges  = errors.New("error while fetching airway bill charges")
	ErrUpsertAWBCharges = errors.New("error while upserting airway bill charges")
	ErrDeleteAWBCharges = errors.New("error while deleting airway bill charges")

	ErrFetchAWBDocs  = errors.New("error while fetching airway bill docs")
	ErrUpsertAWBDocs = errors.New("error while upserting airway bill docs")
	ErrDeleteAWBDocs = errors.New("error while deleting airway bill docs")

	ErrFetchAWBHouse  = errors.New("error while fetching airway bill house")
	ErrUpsertAWBHouse = errors.New("error while upserting airway bill house")
	ErrDeleteAWBHouse = errors.New("error while deleting airway bill house")

	ErrFetchAWBMaster  = errors.New("error while fetching airway bill master")
	ErrUpsertAWBMaster = errors.New("error while upserting airway bill master")
	ErrDeleteAWBMaster = errors.New("error while deleting airway bill master")

	ErrFetchAWBLabels  = errors.New("error while fetching airway bill labels")
	ErrUpsertAWBLabels = errors.New("error while upserting airway bill labels")
	ErrDeleteAWBLabels = errors.New("error while deleting airway bill labels")

	ErrFetchAWBManifest  = errors.New("error while fetching airway bill manifest")
	ErrUpsertAWBManifest = errors.New("error while upserting airway bill manifest")
	ErrDeleteAWBManifest = errors.New("error while deleting airway bill manifest")

	ErrFetchCartingDetails  = errors.New("error while fetching carting details")
	ErrUpsertCartingDetails = errors.New("error while upserting carting details")
	ErrDeleteCartingDetails = errors.New("error while deleting carting details")

	ErrFetchLclBooking  = errors.New("error while fetching lcl booking ")
	ErrUpsertLclBooking = errors.New("error while upserting lcl booking")
	ErrDeleteLclBooking = errors.New("error while deleting lcl booking")

	ErrFetchLclDocsExportCommon  = errors.New("error while fetching common details")
	ErrUpsertLclDocsExportCommon = errors.New("error while upserting common details")
	ErrDeleteLclDocsExportCommon = errors.New("error while deleting common details")

	ErrFetchLclDocsExportShippingInstruction  = errors.New("error while fetching shipping instruction details")
	ErrUpsertLclDocsExportShippingInstruction = errors.New("error while upserting shipping instruction details")
	ErrDeleteLclDocsExportShippingInstruction = errors.New("error while deleting shipping instruction details")

	ErrFetchLclDocsExportVgm  = errors.New("error while fetching  vgm details")
	ErrUpsertLclDocsExportVgm = errors.New("error while upserting vgm details")
	ErrDeleteLclDocsExportVgm = errors.New("error while deleting vgm details")

	ErrLclDocsExportNoConsolId  = errors.New("error while upserting export doc details - consol id is nil! you should always send a consol id")
	ErrLclDocsExportNoBookingId = errors.New("error while upserting export doc details - booking id is nil! you should always send a booking id")

	ErrFetchLcl  = errors.New("error while fetching lcl ")
	ErrUpsertLcl = errors.New("error while upserting lcl")
	ErrDeleteLcl = errors.New("error while deleting lcl")

	ErrFetchLclManifestTask  = errors.New("error while fetching lcl manifest task")
	ErrUpsertLclManifestTask = errors.New("error while upserting lcl manifest task")
	ErrDeleteLclManifestTask = errors.New("error while deleting lcl manifest task")

	ErrFetchLclTsaTask  = errors.New("error while fetching lcl tsa task")
	ErrUpsertLclTsaTask = errors.New("error while upserting lcl tsa task")
	ErrDeleteLclTsaTask = errors.New("error while deleting lcl tsa task")

	ErrFetchADDocs     = errors.New("error while fetching arrival delivery docs")
	ErrUpsertADDocs    = errors.New("error while upserting arrival delivery docs")
	ErrDeleteADDocs    = errors.New("error while deleting arrival delivery docs")
	ErrFetchADDocsDoBl = errors.New("error while accessing addcos DO table")

	ErrFetchShippingDetails  = errors.New("error while fetching Shipping Details")
	ErrUpsertShippingDetails = errors.New("error while upserting Shipping Details")
	ErrDeleteShippingDetails = errors.New("error while deleting Shipping Details")

	ErrFetchStuffingDetails                = errors.New("error while fetching stuffing Details")
	ErrUpsertStuffingDetails               = errors.New("error while upserting stuffing Details")
	ErrDeleteStuffingDetails               = errors.New("error while deleting stuffing Details")
	ErrFetchLclDocsExportPrealertManifest  = errors.New("error while fetching prealert manifest details")
	ErrUpsertLclDocsExportPrealertManifest = errors.New("error while upserting prealert manifest details")
	ErrDeleteLclDocsExportPrealertManifest = errors.New("error while deleting prealert manifest details")

	ErrFetchContainers  = errors.New("error while fetching container details")
	ErrUpsertContainers = errors.New("error while upserting container details")
	ErrDeleteContainers = errors.New("error while deleting container details")

	ErrFetchQuotes = errors.New("error while fetching quotes")
)
