import sys

import healpy as hp
import numpy as np
import matplotlib.pyplot as plt
import random

NSIDE = None
ITERATIONS = 1

def main():
    global NSIDE 
    NSIDE = 2**int(sys.argv[2])
    global ITERATIONS
    ITERATIONS = int(sys.argv[3])

    if sys.argv[1] == "ang2pix":
        test_ang2pix()
    if sys.argv[1] == "conesearch":
        test_conesearch()
    if sys.argv[1] == "conesearch2":
        test_conesearch2()
    if sys.argv[1] == "nside":
        nside_to_resolution(int(sys.argv[2]))
    if sys.argv[1] == "plot-resolutions":
        plot_resolutions()
    if sys.argv[1] == "conesearch-single":
        test_conesearch_single_op()


def test_ang2pix():
    results = []
    for _ in range(0, ITERATIONS):
        for ra in range(0, 360):
            for dec in range(-90, 90):
                theta, phi = radec_to_thetaphi(ra, dec)
                results.append(hp.ang2pix(NSIDE, theta, phi, nest=True))

    print("results:", len(results))
    print(results[:10])

def test_conesearch():
    radius = np.deg2rad(10 / 3600.0)
    results = []
    for _ in range(0, ITERATIONS):
        for ra in range(0, 360):
            for dec in range(-90, 90):
                vec = hp.ang2vec(ra, dec, lonlat=True)
                ipix = hp.query_disc(NSIDE, vec, radius, inclusive=True, nest=True)
                results.append(ipix)

    print("results:", len(results))
    print(results[0])

def test_conesearch2():
    radius = np.deg2rad(10 / 3600.0)
    for _ in range(0, ITERATIONS):
        for ra in range(0, 360):
            for dec in range(-90, 90):
                vec = hp.ang2vec(ra, dec, lonlat=True)
                hp.query_disc(NSIDE, vec, radius, inclusive=True, nest=True)

def test_conesearch_single_op():
    radius = np.deg2rad(10 / 3600.0)
    for _ in range(ITERATIONS):
        ra = random.randint(0, 360)
        dec = random.randint(0, 90)
        vec = hp.ang2vec(ra, dec, lonlat=True)
        hp.query_disc(NSIDE, vec, radius, inclusive=True, nest=True)


def radec_to_thetaphi(ra, dec):
    """
    Convert RA and Dec in degrees to HEALPix theta and phi in radians.
    
    Parameters:
    ra (float): Right Ascension in degrees
    dec (float): Declination in degrees
    
    Returns:
    tuple: (theta, phi) in radians
    """
    theta = np.radians(90 - dec)
    phi = np.radians(ra)
    return theta, phi

def nside_to_resolution(nside):
    value = 2**int(nside)
    print(
        "Approximate resolution at NSIDE {} is {:.2} arcsec".format(
            nside, hp.nside2resol(value, arcmin=True) / 60
        )
    )

def plot_resolutions():
    resolutions = []
    nsides = []
    for nside in range(12, 18):
        nsides.append(nside)
        value = 2**int(nside)
        resolution = hp.nside2resol(value, arcmin=True) / 60
        resolutions.append(resolution)
    plt.plot(nsides, resolutions)
    plt.xlabel("NSIDE")
    plt.ylabel("Resolution (arcsec)")
    plt.show()

if __name__ == "__main__":
    main()
