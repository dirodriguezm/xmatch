API for astronomical cross-matching and catalog queries.

## Overview

CrossWave provides fast cone search and metadata retrieval across multiple astronomical catalogs. The service is optimized for high-throughput queries using HEALPix spatial indexing.

### Key Features

- **Cone Search**: Find objects within a radius of given celestial coordinates
- **Bulk Operations**: Process multiple queries in a single request
- **Metadata Retrieval**: Get detailed catalog information for specific objects
- **Lightcurves**: Retrieve time-series photometry data

## Coordinates

All coordinates use the **J2000 equatorial coordinate system**:

| Parameter | Range | Unit | Description |
|-----------|-------|------|-------------|
| `ra` | 0 to 360 | degrees | Right Ascension |
| `dec` | -90 to +90 | degrees | Declination |
| `radius` | > 0 | degrees | Search radius |

## Available Catalogs

| Catalog | Description |
|---------|-------------|
| `all` | Search across all available catalogs |
| `allwise` | (WISE All-Sky)[https://irsa.ipac.caltech.edu/data/WISE/docs/release/All-Sky/expsup/sec2_2.html] Data Release |
| `gaia` | (Gaia)[https://www.cosmos.esa.int/web/gaia/release] Data Release|
| `erosita` | (eROSITA)[https://erosita.mpe.mpg.de/] (WIP) Data Release|


## Response Codes

| Code | Description |
|------|-------------|
| `200` | Success - results found |
| `204` | Success - no results found |
| `400` | Bad Request - invalid parameters |
| `500` | Internal Server Error |

## Example Usage

### Single Cone Search

Search for objects within 0.01 degrees of RA=180.5, Dec=-45.0:

```
GET /v1/conesearch?ra=180.5&dec=-45.0&radius=0.01&catalog=allwise
```

### Bulk Cone Search

Search multiple positions in a single request:

```json
POST /v1/bulk-conesearch
Content-Type: application/json

{
  "ra": [180.5, 181.2, 182.0],
  "dec": [-45.0, -45.5, -46.0],
  "radius": 0.01,
  "catalog": "allwise",
  "nneighbor": 1
}
```

### Get Metadata by ID

Retrieve detailed catalog data for a specific object:

```
GET /v1/metadata?id=J120000.00-450000.0&catalog=allwise
```

### Get Lightcurve

Retrieve time-series photometry for coordinates:

```
GET /v1/lightcurve?ra=180.5&dec=-45.0&radius=0.01
```

## Error Handling

All validation errors return a JSON object with the following structure:

```json
{
  "field": "ra",
  "reason": "value out of range",
  "value": "400.0"
}
```
