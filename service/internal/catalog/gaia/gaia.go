package gaia

import (
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type InputSchema struct {
	SolutionID                  int64   `parquet:"name=solution_id, type=INT64"`
	Designation                 string  `parquet:"name=designation, type=BYTE_ARRAY"`
	SourceID                    int64   `parquet:"name=source_id, type=INT64"`
	RandomIndex                 int64   `parquet:"name=random_index, type=INT64"`
	RefEpoch                    float64 `parquet:"name=ref_epoch, type=DOUBLE"`
	RA                          float64 `parquet:"name=ra, type=DOUBLE"`
	RAError                     float32 `parquet:"name=ra_error, type=FLOAT"`
	Dec                         float64 `parquet:"name=dec, type=DOUBLE"`
	DecError                    float32 `parquet:"name=dec_error, type=FLOAT"`
	Parallax                    float64 `parquet:"name=parallax, type=DOUBLE"`
	ParallaxError               float32 `parquet:"name=parallax_error, type=FLOAT"`
	ParallaxOverError           float32 `parquet:"name=parallax_over_error, type=FLOAT"`
	PM                          float32 `parquet:"name=pm, type=FLOAT"`
	PMRA                        float64 `parquet:"name=pmra, type=DOUBLE"`
	PMRAError                   float32 `parquet:"name=pmra_error, type=FLOAT"`
	PMDec                       float64 `parquet:"name=pmdec, type=DOUBLE"`
	PMDecError                  float32 `parquet:"name=pmdec_error, type=FLOAT"`
	RaDecCorr                   float32 `parquet:"name=ra_dec_corr, type=FLOAT"`
	RaParallaxCorr              float32 `parquet:"name=ra_parallax_corr, type=FLOAT"`
	RaPmraCorr                  float32 `parquet:"name=ra_pmra_corr, type=FLOAT"`
	RaPmdecCorr                 float32 `parquet:"name=ra_pmdec_corr, type=FLOAT"`
	DecParallaxCorr             float32 `parquet:"name=dec_parallax_corr, type=FLOAT"`
	DecPmraCorr                 float32 `parquet:"name=dec_pmra_corr, type=FLOAT"`
	DecPmdecCorr                float32 `parquet:"name=dec_pmdec_corr, type=FLOAT"`
	ParallaxPmraCorr            float32 `parquet:"name=parallax_pmra_corr, type=FLOAT"`
	ParallaxPmdecCorr           float32 `parquet:"name=parallax_pmdec_corr, type=FLOAT"`
	PmraPmdecCorr               float32 `parquet:"name=pmra_pmdec_corr, type=FLOAT"`
	AstrometricNObsAl           int16   `parquet:"name=astrometric_n_obs_al, type=INT32"`
	AstrometricNObsAc           int16   `parquet:"name=astrometric_n_obs_ac, type=INT32"`
	AstrometricNGoodObsAl       int16   `parquet:"name=astrometric_n_good_obs_al, type=INT32"`
	AstrometricNBadObsAl        int16   `parquet:"name=astrometric_n_bad_obs_al, type=INT32"`
	AstrometricGofAl            float32 `parquet:"name=astrometric_gof_al, type=FLOAT"`
	AstrometricChi2Al           float32 `parquet:"name=astrometric_chi2_al, type=FLOAT"`
	AstrometricExcessNoise      float32 `parquet:"name=astrometric_excess_noise, type=FLOAT"`
	AstrometricExcessNoiseSig   float32 `parquet:"name=astrometric_excess_noise_sig, type=FLOAT"`
	AstrometricParamsSolved     int8    `parquet:"name=astrometric_params_solved, type=INT32"`
	AstrometricPrimaryFlag      bool    `parquet:"name=astrometric_primary_flag, type=BOOLEAN"`
	NuEffUsedInAstrometry       float32 `parquet:"name=nu_eff_used_in_astrometry, type=FLOAT"`
	Pseudocolour                float32 `parquet:"name=pseudocolour, type=FLOAT"`
	PseudocolourError           float32 `parquet:"name=pseudocolour_error, type=FLOAT"`
	RaPseudocolourCorr          float32 `parquet:"name=ra_pseudocolour_corr, type=FLOAT"`
	DecPseudocolourCorr         float32 `parquet:"name=dec_pseudocolour_corr, type=FLOAT"`
	ParallaxPseudocolourCorr    float32 `parquet:"name=parallax_pseudocolour_corr, type=FLOAT"`
	PmraPseudocolourCorr        float32 `parquet:"name=pmra_pseudocolour_corr, type=FLOAT"`
	PmdecPseudocolourCorr       float32 `parquet:"name=pmdec_pseudocolour_corr, type=FLOAT"`
	AstrometricMatchedTransits  int16   `parquet:"name=astrometric_matched_transits, type=INT32"`
	VisibilityPeriodsUsed       int16   `parquet:"name=visibility_periods_used, type=INT32"`
	AstrometricSigma5dMax       float32 `parquet:"name=astrometric_sigma5d_max, type=FLOAT"`
	MatchedTransits             int16   `parquet:"name=matched_transits, type=INT32"`
	NewMatchedTransits          int16   `parquet:"name=new_matched_transits, type=INT32"`
	MatchedTransitsRemoved      int16   `parquet:"name=matched_transits_removed, type=INT32"`
	IpdGofHarmonicAmplitude     float32 `parquet:"name=ipd_gof_harmonic_amplitude, type=FLOAT"`
	IpdGofHarmonicPhase         float32 `parquet:"name=ipd_gof_harmonic_phase, type=FLOAT"`
	IpdFracMultiPeak            int8    `parquet:"name=ipd_frac_multi_peak, type=INT32"`
	IpdFracOddWin               int8    `parquet:"name=ipd_frac_odd_win, type=INT32"`
	Ruwe                        float32 `parquet:"name=ruwe, type=FLOAT"`
	ScanDirectionStrengthK1     float32 `parquet:"name=scan_direction_strength_k1, type=FLOAT"`
	ScanDirectionStrengthK2     float32 `parquet:"name=scan_direction_strength_k2, type=FLOAT"`
	ScanDirectionStrengthK3     float32 `parquet:"name=scan_direction_strength_k3, type=FLOAT"`
	ScanDirectionStrengthK4     float32 `parquet:"name=scan_direction_strength_k4, type=FLOAT"`
	ScanDirectionMeanK1         float32 `parquet:"name=scan_direction_mean_k1, type=FLOAT"`
	ScanDirectionMeanK2         float32 `parquet:"name=scan_direction_mean_k2, type=FLOAT"`
	ScanDirectionMeanK3         float32 `parquet:"name=scan_direction_mean_k3, type=FLOAT"`
	ScanDirectionMeanK4         float32 `parquet:"name=scan_direction_mean_k4, type=FLOAT"`
	DuplicatedSource            bool    `parquet:"name=duplicated_source, type=BOOLEAN"`
	PhotGNObs                   int16   `parquet:"name=phot_g_n_obs, type=INT32"`
	PhotGMeanFlux               float64 `parquet:"name=phot_g_mean_flux, type=DOUBLE"`
	PhotGMeanFluxError          float32 `parquet:"name=phot_g_mean_flux_error, type=FLOAT"`
	PhotGMeanFluxOverError      float32 `parquet:"name=phot_g_mean_flux_over_error, type=FLOAT"`
	PhotGMeanMag                float32 `parquet:"name=phot_g_mean_mag, type=FLOAT"`
	PhotBpNObs                  int16   `parquet:"name=phot_bp_n_obs, type=INT32"`
	PhotBpMeanFlux              float64 `parquet:"name=phot_bp_mean_flux, type=DOUBLE"`
	PhotBpMeanFluxError         float32 `parquet:"name=phot_bp_mean_flux_error, type=FLOAT"`
	PhotBpMeanFluxOverError     float32 `parquet:"name=phot_bp_mean_flux_over_error, type=FLOAT"`
	PhotBpMeanMag               float32 `parquet:"name=phot_bp_mean_mag, type=FLOAT"`
	PhotRpNObs                  int16   `parquet:"name=phot_rp_n_obs, type=INT32"`
	PhotRpMeanFlux              float64 `parquet:"name=phot_rp_mean_flux, type=DOUBLE"`
	PhotRpMeanFluxError         float32 `parquet:"name=phot_rp_mean_flux_error, type=FLOAT"`
	PhotRpMeanFluxOverError     float32 `parquet:"name=phot_rp_mean_flux_over_error, type=FLOAT"`
	PhotRpMeanMag               float32 `parquet:"name=phot_rp_mean_mag, type=FLOAT"`
	PhotBpRpExcessFactor        float32 `parquet:"name=phot_bp_rp_excess_factor, type=FLOAT"`
	PhotBpNContaminatedTransits int16   `parquet:"name=phot_bp_n_contaminated_transits, type=INT32"`
	PhotBpNBlendedTransits      int16   `parquet:"name=phot_bp_n_blended_transits, type=INT32"`
	PhotRpNContaminatedTransits int16   `parquet:"name=phot_rp_n_contaminated_transits, type=INT32"`
	PhotRpNBlendedTransits      int16   `parquet:"name=phot_rp_n_blended_transits, type=INT32"`
	PhotProcMode                int8    `parquet:"name=phot_proc_mode, type=INT32"`
	BpRp                        float32 `parquet:"name=bp_rp, type=FLOAT"`
	BpG                         float32 `parquet:"name=bp_g, type=FLOAT"`
	GRp                         float32 `parquet:"name=g_rp, type=FLOAT"`
	RadialVelocity              float32 `parquet:"name=radial_velocity, type=FLOAT"`
	RadialVelocityError         float32 `parquet:"name=radial_velocity_error, type=FLOAT"`
	RvMethodUsed                int8    `parquet:"name=rv_method_used, type=INT32"`
	RvNbTransits                int16   `parquet:"name=rv_nb_transits, type=INT32"`
	RvNbDeblendedTransits       int16   `parquet:"name=rv_nb_deblended_transits, type=INT32"`
	RvVisibilityPeriodsUsed     int16   `parquet:"name=rv_visibility_periods_used, type=INT32"`
	RvExpectedSigToNoise        float32 `parquet:"name=rv_expected_sig_to_noise, type=FLOAT"`
	RvRenormalisedGof           float32 `parquet:"name=rv_renormalised_gof, type=FLOAT"`
	RvChisqPvalue               float32 `parquet:"name=rv_chisq_pvalue, type=FLOAT"`
	RvTimeDuration              float32 `parquet:"name=rv_time_duration, type=FLOAT"`
	RvAmplitudeRobust           float32 `parquet:"name=rv_amplitude_robust, type=FLOAT"`
	RvTemplateTeff              float32 `parquet:"name=rv_template_teff, type=FLOAT"`
	RvTemplateLogg              float32 `parquet:"name=rv_template_logg, type=FLOAT"`
	RvTemplateFeH               float32 `parquet:"name=rv_template_fe_h, type=FLOAT"`
	RvAtmParamOrigin            int16   `parquet:"name=rv_atm_param_origin, type=INT32"`
	Vbroad                      float32 `parquet:"name=vbroad, type=FLOAT"`
	VbroadError                 float32 `parquet:"name=vbroad_error, type=FLOAT"`
	VbroadNbTransits            int16   `parquet:"name=vbroad_nb_transits, type=INT32"`
	GrvsMag                     float32 `parquet:"name=grvs_mag, type=FLOAT"`
	GrvsMagError                float32 `parquet:"name=grvs_mag_error, type=FLOAT"`
	GrvsMagNbTransits           int16   `parquet:"name=grvs_mag_nb_transits, type=INT32"`
	RvsSpecSigToNoise           float32 `parquet:"name=rvs_spec_sig_to_noise, type=FLOAT"`
	PhotVariableFlag            string  `parquet:"name=phot_variable_flag, type=BYTE_ARRAY"`
	L                           float64 `parquet:"name=l, type=DOUBLE"`
	B                           float64 `parquet:"name=b, type=DOUBLE"`
	EclLon                      float64 `parquet:"name=ecl_lon, type=DOUBLE"`
	EclLat                      float64 `parquet:"name=ecl_lat, type=DOUBLE"`
	InQsoCandidates             bool    `parquet:"name=in_qso_candidates, type=BOOLEAN"`
	InGalaxyCandidates          bool    `parquet:"name=in_galaxy_candidates, type=BOOLEAN"`
	NonSingleStar               int16   `parquet:"name=non_single_star, type=INT32"`
	HasXpContinuous             bool    `parquet:"name=has_xp_continuous, type=BOOLEAN"`
	HasXpSampled                bool    `parquet:"name=has_xp_sampled, type=BOOLEAN"`
	HasRvs                      bool    `parquet:"name=has_rvs, type=BOOLEAN"`
	HasEpochPhotometry          bool    `parquet:"name=has_epoch_photometry, type=BOOLEAN"`
	HasEpochRv                  bool    `parquet:"name=has_epoch_rv, type=BOOLEAN"`
	HasMcmcGspphot              bool    `parquet:"name=has_mcmc_gspphot, type=BOOLEAN"`
	HasMcmcMsc                  bool    `parquet:"name=has_mcmc_msc, type=BOOLEAN"`
	InAndromedaSurvey           bool    `parquet:"name=in_andromeda_survey, type=BOOLEAN"`
	ClassprobDscCombmodQuasar   float32 `parquet:"name=classprob_dsc_combmod_quasar, type=FLOAT"`
	ClassprobDscCombmodGalaxy   float32 `parquet:"name=classprob_dsc_combmod_galaxy, type=FLOAT"`
	ClassprobDscCombmodStar     float32 `parquet:"name=classprob_dsc_combmod_star, type=FLOAT"`
	TeffGspphot                 float32 `parquet:"name=teff_gspphot, type=FLOAT"`
	TeffGspphotLower            float32 `parquet:"name=teff_gspphot_lower, type=FLOAT"`
	TeffGspphotUpper            float32 `parquet:"name=teff_gspphot_upper, type=FLOAT"`
	LoggGspphot                 float32 `parquet:"name=logg_gspphot, type=FLOAT"`
	LoggGspphotLower            float32 `parquet:"name=logg_gspphot_lower, type=FLOAT"`
	LoggGspphotUpper            float32 `parquet:"name=logg_gspphot_upper, type=FLOAT"`
	MhGspphot                   float32 `parquet:"name=mh_gspphot, type=FLOAT"`
	MhGspphotLower              float32 `parquet:"name=mh_gspphot_lower, type=FLOAT"`
	MhGspphotUpper              float32 `parquet:"name=mh_gspphot_upper, type=FLOAT"`
	DistanceGspphot             float32 `parquet:"name=distance_gspphot, type=FLOAT"`
	DistanceGspphotLower        float32 `parquet:"name=distance_gspphot_lower, type=FLOAT"`
	DistanceGspphotUpper        float32 `parquet:"name=distance_gspphot_upper, type=FLOAT"`
	AzeroGspphot                float32 `parquet:"name=azero_gspphot, type=FLOAT"`
	AzeroGspphotLower           float32 `parquet:"name=azero_gspphot_lower, type=FLOAT"`
	AzeroGspphotUpper           float32 `parquet:"name=azero_gspphot_upper, type=FLOAT"`
	AgGspphot                   float32 `parquet:"name=ag_gspphot, type=FLOAT"`
	AgGspphotLower              float32 `parquet:"name=ag_gspphot_lower, type=FLOAT"`
	AgGspphotUpper              float32 `parquet:"name=ag_gspphot_upper, type=FLOAT"`
	EbpminrpGspphot             float32 `parquet:"name=ebpminrp_gspphot, type=FLOAT"`
	EbpminrpGspphotLower        float32 `parquet:"name=ebpminrp_gspphot_lower, type=FLOAT"`
	EbpminrpGspphotUpper        float32 `parquet:"name=ebpminrp_gspphot_upper, type=FLOAT"`
	LibnameGspphot              string  `parquet:"name=libname_gspphot, type=BYTE_ARRAY"`
}

func boolToInt64(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func (r InputSchema) GetCoordinates() (float64, float64) {
	return r.RA, r.Dec
}

func (r InputSchema) MastercatID() string {
	return r.Designation
}

func (schema InputSchema) ConvertToMetadata() (repository.Gaia, error) {
	return repository.Gaia{
		ID:                        schema.Designation,
		SourceID:                  new(schema.SourceID),
		RaError:                   new(float64(schema.RAError)),
		DecError:                  new(float64(schema.DecError)),
		Parallax:                  new(schema.Parallax),
		ParallaxError:             new(float64(schema.ParallaxError)),
		Pm:                        new(float64(schema.PM)),
		Pmra:                      new(schema.PMRA),
		PmraError:                 new(float64(schema.PMRAError)),
		Pmdec:                     new(schema.PMDec),
		PmdecError:                new(float64(schema.PMDecError)),
		AstrometricExcessNoise:    new(float64(schema.AstrometricExcessNoise)),
		AstrometricExcessNoiseSig: new(float64(schema.AstrometricExcessNoiseSig)),
		Ruwe:                      new(float64(schema.Ruwe)),
		PhotGNObs:                 new(int64(schema.PhotGNObs)),
		PhotGMeanFlux:             new(schema.PhotGMeanFlux),
		PhotGMeanFluxError:        new(float64(schema.PhotGMeanFluxError)),
		PhotGMeanFluxOverError:    new(float64(schema.PhotGMeanFluxOverError)),
		PhotGMeanMag:              new(float64(schema.PhotGMeanMag)),
		PhotBpNObs:                new(int64(schema.PhotBpNObs)),
		PhotBpMeanFlux:            new(schema.PhotBpMeanFlux),
		PhotBpMeanFluxError:       new(float64(schema.PhotBpMeanFluxError)),
		PhotBpMeanFluxOverError:   new(float64(schema.PhotBpMeanFluxOverError)),
		PhotBpMeanMag:             new(float64(schema.PhotBpMeanMag)),
		PhotRpNObs:                new(int64(schema.PhotRpNObs)),
		PhotRpMeanFlux:            new(schema.PhotRpMeanFlux),
		PhotRpMeanFluxError:       new(float64(schema.PhotRpMeanFluxError)),
		PhotRpMeanFluxOverError:   new(float64(schema.PhotRpMeanFluxOverError)),
		PhotRpMeanMag:             new(float64(schema.PhotRpMeanMag)),
		PhotBpRpExcessFactor:      new(float64(schema.PhotBpRpExcessFactor)),
		PhotProcMode:              new(int64(schema.PhotProcMode)),
		BpRp:                      new(float64(schema.BpRp)),
		BpG:                       new(float64(schema.BpG)),
		GRp:                       new(float64(schema.GRp)),
		RadialVelocity:            new(float64(schema.RadialVelocity)),
		RadialVelocityError:       new(float64(schema.RadialVelocityError)),
		RvMethodUsed:              new(int64(schema.RvMethodUsed)),
		PhotVariableFlag:          new(schema.PhotVariableFlag),
		InQsoCandidates:           new(boolToInt64(schema.InQsoCandidates)),
		InGalaxyCandidates:        new(boolToInt64(schema.InGalaxyCandidates)),
		NonSingleStar:             new(int64(schema.NonSingleStar)),
		HasEpochPhotometry:        new(boolToInt64(schema.HasEpochPhotometry)),
		ClassprobDscCombmodQuasar: new(float64(schema.ClassprobDscCombmodQuasar)),
		ClassprobDscCombmodGalaxy: new(float64(schema.ClassprobDscCombmodGalaxy)),
		ClassprobDscCombmodStar:   new(float64(schema.ClassprobDscCombmodStar)),
		TeffGspphot:               new(float64(schema.TeffGspphot)),
		TeffGspphotLower:          new(float64(schema.TeffGspphotLower)),
		TeffGspphotUpper:          new(float64(schema.TeffGspphotUpper)),
		LoggGspphot:               new(float64(schema.LoggGspphot)),
		LoggGspphotLower:          new(float64(schema.LoggGspphotLower)),
		LoggGspphotUpper:          new(float64(schema.LoggGspphotUpper)),
		MhGspphot:                 new(float64(schema.MhGspphot)),
		MhGspphotLower:            new(float64(schema.MhGspphotLower)),
		MhGspphotUpper:            new(float64(schema.MhGspphotUpper)),
		DistanceGspphot:           new(float64(schema.DistanceGspphot)),
		DistanceGspphotLower:      new(float64(schema.DistanceGspphotLower)),
		DistanceGspphotUpper:      new(float64(schema.DistanceGspphotUpper)),
	}, nil
}
