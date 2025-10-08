package repository

import (
	"context"
	"database/sql"
)

type GaiaInputSchema struct {
	SolutionID                  int64   `json:"solution_id" parquet:"name=solution_id, type=INT64"`
	Designation                 string  `json:"designation" parquet:"name=designation, type=BYTE_ARRAY"`
	SourceID                    int64   `json:"source_id" parquet:"name=source_id, type=INT64"`
	RandomIndex                 int64   `json:"random_index" parquet:"name=random_index, type=INT64"`
	RefEpoch                    float64 `json:"ref_epoch" parquet:"name=ref_epoch, type=DOUBLE"`
	RA                          float64 `json:"ra" parquet:"name=ra, type=DOUBLE"`
	RAError                     float32 `json:"ra_error" parquet:"name=ra_error, type=FLOAT"`
	Dec                         float64 `json:"dec" parquet:"name=dec, type=DOUBLE"`
	DecError                    float32 `json:"dec_error" parquet:"name=dec_error, type=FLOAT"`
	Parallax                    float64 `json:"parallax" parquet:"name=parallax, type=DOUBLE"`
	ParallaxError               float32 `json:"parallax_error" parquet:"name=parallax_error, type=FLOAT"`
	ParallaxOverError           float32 `json:"parallax_over_error" parquet:"name=parallax_over_error, type=FLOAT"`
	PM                          float32 `json:"pm" parquet:"name=pm, type=FLOAT"`
	PMRA                        float64 `json:"pmra" parquet:"name=pmra, type=DOUBLE"`
	PMRAError                   float32 `json:"pmra_error" parquet:"name=pmra_error, type=FLOAT"`
	PMDec                       float64 `json:"pmdec" parquet:"name=pmdec, type=DOUBLE"`
	PMDecError                  float32 `json:"pmdec_error" parquet:"name=pmdec_error, type=FLOAT"`
	RaDecCorr                   float32 `json:"ra_dec_corr" parquet:"name=ra_dec_corr, type=FLOAT"`
	RaParallaxCorr              float32 `json:"ra_parallax_corr" parquet:"name=ra_parallax_corr, type=FLOAT"`
	RaPmraCorr                  float32 `json:"ra_pmra_corr" parquet:"name=ra_pmra_corr, type=FLOAT"`
	RaPmdecCorr                 float32 `json:"ra_pmdec_corr" parquet:"name=ra_pmdec_corr, type=FLOAT"`
	DecParallaxCorr             float32 `json:"dec_parallax_corr" parquet:"name=dec_parallax_corr, type=FLOAT"`
	DecPmraCorr                 float32 `json:"dec_pmra_corr" parquet:"name=dec_pmra_corr, type=FLOAT"`
	DecPmdecCorr                float32 `json:"dec_pmdec_corr" parquet:"name=dec_pmdec_corr, type=FLOAT"`
	ParallaxPmraCorr            float32 `json:"parallax_pmra_corr" parquet:"name=parallax_pmra_corr, type=FLOAT"`
	ParallaxPmdecCorr           float32 `json:"parallax_pmdec_corr" parquet:"name=parallax_pmdec_corr, type=FLOAT"`
	PmraPmdecCorr               float32 `json:"pmra_pmdec_corr" parquet:"name=pmra_pmdec_corr, type=FLOAT"`
	AstrometricNObsAl           int16   `json:"astrometric_n_obs_al" parquet:"name=astrometric_n_obs_al, type=INT32"`
	AstrometricNObsAc           int16   `json:"astrometric_n_obs_ac" parquet:"name=astrometric_n_obs_ac, type=INT32"`
	AstrometricNGoodObsAl       int16   `json:"astrometric_n_good_obs_al" parquet:"name=astrometric_n_good_obs_al, type=INT32"`
	AstrometricNBadObsAl        int16   `json:"astrometric_n_bad_obs_al" parquet:"name=astrometric_n_bad_obs_al, type=INT32"`
	AstrometricGofAl            float32 `json:"astrometric_gof_al" parquet:"name=astrometric_gof_al, type=FLOAT"`
	AstrometricChi2Al           float32 `json:"astrometric_chi2_al" parquet:"name=astrometric_chi2_al, type=FLOAT"`
	AstrometricExcessNoise      float32 `json:"astrometric_excess_noise" parquet:"name=astrometric_excess_noise, type=FLOAT"`
	AstrometricExcessNoiseSig   float32 `json:"astrometric_excess_noise_sig" parquet:"name=astrometric_excess_noise_sig, type=FLOAT"`
	AstrometricParamsSolved     int8    `json:"astrometric_params_solved" parquet:"name=astrometric_params_solved, type=INT32"`
	AstrometricPrimaryFlag      bool    `json:"astrometric_primary_flag" parquet:"name=astrometric_primary_flag, type=BOOLEAN"`
	NuEffUsedInAstrometry       float32 `json:"nu_eff_used_in_astrometry" parquet:"name=nu_eff_used_in_astrometry, type=FLOAT"`
	Pseudocolour                float32 `json:"pseudocolour" parquet:"name=pseudocolour, type=FLOAT"`
	PseudocolourError           float32 `json:"pseudocolour_error" parquet:"name=pseudocolour_error, type=FLOAT"`
	RaPseudocolourCorr          float32 `json:"ra_pseudocolour_corr" parquet:"name=ra_pseudocolour_corr, type=FLOAT"`
	DecPseudocolourCorr         float32 `json:"dec_pseudocolour_corr" parquet:"name=dec_pseudocolour_corr, type=FLOAT"`
	ParallaxPseudocolourCorr    float32 `json:"parallax_pseudocolour_corr" parquet:"name=parallax_pseudocolour_corr, type=FLOAT"`
	PmraPseudocolourCorr        float32 `json:"pmra_pseudocolour_corr" parquet:"name=pmra_pseudocolour_corr, type=FLOAT"`
	PmdecPseudocolourCorr       float32 `json:"pmdec_pseudocolour_corr" parquet:"name=pmdec_pseudocolour_corr, type=FLOAT"`
	AstrometricMatchedTransits  int16   `json:"astrometric_matched_transits" parquet:"name=astrometric_matched_transits, type=INT32"`
	VisibilityPeriodsUsed       int16   `json:"visibility_periods_used" parquet:"name=visibility_periods_used, type=INT32"`
	AstrometricSigma5dMax       float32 `json:"astrometric_sigma5d_max" parquet:"name=astrometric_sigma5d_max, type=FLOAT"`
	MatchedTransits             int16   `json:"matched_transits" parquet:"name=matched_transits, type=INT32"`
	NewMatchedTransits          int16   `json:"new_matched_transits" parquet:"name=new_matched_transits, type=INT32"`
	MatchedTransitsRemoved      int16   `json:"matched_transits_removed" parquet:"name=matched_transits_removed, type=INT32"`
	IpdGofHarmonicAmplitude     float32 `json:"ipd_gof_harmonic_amplitude" parquet:"name=ipd_gof_harmonic_amplitude, type=FLOAT"`
	IpdGofHarmonicPhase         float32 `json:"ipd_gof_harmonic_phase" parquet:"name=ipd_gof_harmonic_phase, type=FLOAT"`
	IpdFracMultiPeak            int8    `json:"ipd_frac_multi_peak" parquet:"name=ipd_frac_multi_peak, type=INT32"`
	IpdFracOddWin               int8    `json:"ipd_frac_odd_win" parquet:"name=ipd_frac_odd_win, type=INT32"`
	Ruwe                        float32 `json:"ruwe" parquet:"name=ruwe, type=FLOAT"`
	ScanDirectionStrengthK1     float32 `json:"scan_direction_strength_k1" parquet:"name=scan_direction_strength_k1, type=FLOAT"`
	ScanDirectionStrengthK2     float32 `json:"scan_direction_strength_k2" parquet:"name=scan_direction_strength_k2, type=FLOAT"`
	ScanDirectionStrengthK3     float32 `json:"scan_direction_strength_k3" parquet:"name=scan_direction_strength_k3, type=FLOAT"`
	ScanDirectionStrengthK4     float32 `json:"scan_direction_strength_k4" parquet:"name=scan_direction_strength_k4, type=FLOAT"`
	ScanDirectionMeanK1         float32 `json:"scan_direction_mean_k1" parquet:"name=scan_direction_mean_k1, type=FLOAT"`
	ScanDirectionMeanK2         float32 `json:"scan_direction_mean_k2" parquet:"name=scan_direction_mean_k2, type=FLOAT"`
	ScanDirectionMeanK3         float32 `json:"scan_direction_mean_k3" parquet:"name=scan_direction_mean_k3, type=FLOAT"`
	ScanDirectionMeanK4         float32 `json:"scan_direction_mean_k4" parquet:"name=scan_direction_mean_k4, type=FLOAT"`
	DuplicatedSource            bool    `json:"duplicated_source" parquet:"name=duplicated_source, type=BOOLEAN"`
	PhotGNObs                   int16   `json:"phot_g_n_obs" parquet:"name=phot_g_n_obs, type=INT32"`
	PhotGMeanFlux               float64 `json:"phot_g_mean_flux" parquet:"name=phot_g_mean_flux, type=DOUBLE"`
	PhotGMeanFluxError          float32 `json:"phot_g_mean_flux_error" parquet:"name=phot_g_mean_flux_error, type=FLOAT"`
	PhotGMeanFluxOverError      float32 `json:"phot_g_mean_flux_over_error" parquet:"name=phot_g_mean_flux_over_error, type=FLOAT"`
	PhotGMeanMag                float32 `json:"phot_g_mean_mag" parquet:"name=phot_g_mean_mag, type=FLOAT"`
	PhotBpNObs                  int16   `json:"phot_bp_n_obs" parquet:"name=phot_bp_n_obs, type=INT32"`
	PhotBpMeanFlux              float64 `json:"phot_bp_mean_flux" parquet:"name=phot_bp_mean_flux, type=DOUBLE"`
	PhotBpMeanFluxError         float32 `json:"phot_bp_mean_flux_error" parquet:"name=phot_bp_mean_flux_error, type=FLOAT"`
	PhotBpMeanFluxOverError     float32 `json:"phot_bp_mean_flux_over_error" parquet:"name=phot_bp_mean_flux_over_error, type=FLOAT"`
	PhotBpMeanMag               float32 `json:"phot_bp_mean_mag" parquet:"name=phot_bp_mean_mag, type=FLOAT"`
	PhotRpNObs                  int16   `json:"phot_rp_n_obs" parquet:"name=phot_rp_n_obs, type=INT32"`
	PhotRpMeanFlux              float64 `json:"phot_rp_mean_flux" parquet:"name=phot_rp_mean_flux, type=DOUBLE"`
	PhotRpMeanFluxError         float32 `json:"phot_rp_mean_flux_error" parquet:"name=phot_rp_mean_flux_error, type=FLOAT"`
	PhotRpMeanFluxOverError     float32 `json:"phot_rp_mean_flux_over_error" parquet:"name=phot_rp_mean_flux_over_error, type=FLOAT"`
	PhotRpMeanMag               float32 `json:"phot_rp_mean_mag" parquet:"name=phot_rp_mean_mag, type=FLOAT"`
	PhotBpRpExcessFactor        float32 `json:"phot_bp_rp_excess_factor" parquet:"name=phot_bp_rp_excess_factor, type=FLOAT"`
	PhotBpNContaminatedTransits int16   `json:"phot_bp_n_contaminated_transits" parquet:"name=phot_bp_n_contaminated_transits, type=INT32"`
	PhotBpNBlendedTransits      int16   `json:"phot_bp_n_blended_transits" parquet:"name=phot_bp_n_blended_transits, type=INT32"`
	PhotRpNContaminatedTransits int16   `json:"phot_rp_n_contaminated_transits" parquet:"name=phot_rp_n_contaminated_transits, type=INT32"`
	PhotRpNBlendedTransits      int16   `json:"phot_rp_n_blended_transits" parquet:"name=phot_rp_n_blended_transits, type=INT32"`
	PhotProcMode                int8    `json:"phot_proc_mode" parquet:"name=phot_proc_mode, type=INT32"`
	BpRp                        float32 `json:"bp_rp" parquet:"name=bp_rp, type=FLOAT"`
	BpG                         float32 `json:"bp_g" parquet:"name=bp_g, type=FLOAT"`
	GRp                         float32 `json:"g_rp" parquet:"name=g_rp, type=FLOAT"`
	RadialVelocity              float32 `json:"radial_velocity" parquet:"name=radial_velocity, type=FLOAT"`
	RadialVelocityError         float32 `json:"radial_velocity_error" parquet:"name=radial_velocity_error, type=FLOAT"`
	RvMethodUsed                int8    `json:"rv_method_used" parquet:"name=rv_method_used, type=INT32"`
	RvNbTransits                int16   `json:"rv_nb_transits" parquet:"name=rv_nb_transits, type=INT32"`
	RvNbDeblendedTransits       int16   `json:"rv_nb_deblended_transits" parquet:"name=rv_nb_deblended_transits, type=INT32"`
	RvVisibilityPeriodsUsed     int16   `json:"rv_visibility_periods_used" parquet:"name=rv_visibility_periods_used, type=INT32"`
	RvExpectedSigToNoise        float32 `json:"rv_expected_sig_to_noise" parquet:"name=rv_expected_sig_to_noise, type=FLOAT"`
	RvRenormalisedGof           float32 `json:"rv_renormalised_gof" parquet:"name=rv_renormalised_gof, type=FLOAT"`
	RvChisqPvalue               float32 `json:"rv_chisq_pvalue" parquet:"name=rv_chisq_pvalue, type=FLOAT"`
	RvTimeDuration              float32 `json:"rv_time_duration" parquet:"name=rv_time_duration, type=FLOAT"`
	RvAmplitudeRobust           float32 `json:"rv_amplitude_robust" parquet:"name=rv_amplitude_robust, type=FLOAT"`
	RvTemplateTeff              float32 `json:"rv_template_teff" parquet:"name=rv_template_teff, type=FLOAT"`
	RvTemplateLogg              float32 `json:"rv_template_logg" parquet:"name=rv_template_logg, type=FLOAT"`
	RvTemplateFeH               float32 `json:"rv_template_fe_h" parquet:"name=rv_template_fe_h, type=FLOAT"`
	RvAtmParamOrigin            int16   `json:"rv_atm_param_origin" parquet:"name=rv_atm_param_origin, type=INT32"`
	Vbroad                      float32 `json:"vbroad" parquet:"name=vbroad, type=FLOAT"`
	VbroadError                 float32 `json:"vbroad_error" parquet:"name=vbroad_error, type=FLOAT"`
	VbroadNbTransits            int16   `json:"vbroad_nb_transits" parquet:"name=vbroad_nb_transits, type=INT32"`
	GrvsMag                     float32 `json:"grvs_mag" parquet:"name=grvs_mag, type=FLOAT"`
	GrvsMagError                float32 `json:"grvs_mag_error" parquet:"name=grvs_mag_error, type=FLOAT"`
	GrvsMagNbTransits           int16   `json:"grvs_mag_nb_transits" parquet:"name=grvs_mag_nb_transits, type=INT32"`
	RvsSpecSigToNoise           float32 `json:"rvs_spec_sig_to_noise" parquet:"name=rvs_spec_sig_to_noise, type=FLOAT"`
	PhotVariableFlag            string  `json:"phot_variable_flag" parquet:"name=phot_variable_flag, type=BYTE_ARRAY"`
	L                           float64 `json:"l" parquet:"name=l, type=DOUBLE"`
	B                           float64 `json:"b" parquet:"name=b, type=DOUBLE"`
	EclLon                      float64 `json:"ecl_lon" parquet:"name=ecl_lon, type=DOUBLE"`
	EclLat                      float64 `json:"ecl_lat" parquet:"name=ecl_lat, type=DOUBLE"`
	InQsoCandidates             bool    `json:"in_qso_candidates" parquet:"name=in_qso_candidates, type=BOOLEAN"`
	InGalaxyCandidates          bool    `json:"in_galaxy_candidates" parquet:"name=in_galaxy_candidates, type=BOOLEAN"`
	NonSingleStar               int16   `json:"non_single_star" parquet:"name=non_single_star, type=INT32"`
	HasXpContinuous             bool    `json:"has_xp_continuous" parquet:"name=has_xp_continuous, type=BOOLEAN"`
	HasXpSampled                bool    `json:"has_xp_sampled" parquet:"name=has_xp_sampled, type=BOOLEAN"`
	HasRvs                      bool    `json:"has_rvs" parquet:"name=has_rvs, type=BOOLEAN"`
	HasEpochPhotometry          bool    `json:"has_epoch_photometry" parquet:"name=has_epoch_photometry, type=BOOLEAN"`
	HasEpochRv                  bool    `json:"has_epoch_rv" parquet:"name=has_epoch_rv, type=BOOLEAN"`
	HasMcmcGspphot              bool    `json:"has_mcmc_gspphot" parquet:"name=has_mcmc_gspphot, type=BOOLEAN"`
	HasMcmcMsc                  bool    `json:"has_mcmc_msc" parquet:"name=has_mcmc_msc, type=BOOLEAN"`
	InAndromedaSurvey           bool    `json:"in_andromeda_survey" parquet:"name=in_andromeda_survey, type=BOOLEAN"`
	ClassprobDscCombmodQuasar   float32 `json:"classprob_dsc_combmod_quasar" parquet:"name=classprob_dsc_combmod_quasar, type=FLOAT"`
	ClassprobDscCombmodGalaxy   float32 `json:"classprob_dsc_combmod_galaxy" parquet:"name=classprob_dsc_combmod_galaxy, type=FLOAT"`
	ClassprobDscCombmodStar     float32 `json:"classprob_dsc_combmod_star" parquet:"name=classprob_dsc_combmod_star, type=FLOAT"`
	TeffGspphot                 float32 `json:"teff_gspphot" parquet:"name=teff_gspphot, type=FLOAT"`
	TeffGspphotLower            float32 `json:"teff_gspphot_lower" parquet:"name=teff_gspphot_lower, type=FLOAT"`
	TeffGspphotUpper            float32 `json:"teff_gspphot_upper" parquet:"name=teff_gspphot_upper, type=FLOAT"`
	LoggGspphot                 float32 `json:"logg_gspphot" parquet:"name=logg_gspphot, type=FLOAT"`
	LoggGspphotLower            float32 `json:"logg_gspphot_lower" parquet:"name=logg_gspphot_lower, type=FLOAT"`
	LoggGspphotUpper            float32 `json:"logg_gspphot_upper" parquet:"name=logg_gspphot_upper, type=FLOAT"`
	MhGspphot                   float32 `json:"mh_gspphot" parquet:"name=mh_gspphot, type=FLOAT"`
	MhGspphotLower              float32 `json:"mh_gspphot_lower" parquet:"name=mh_gspphot_lower, type=FLOAT"`
	MhGspphotUpper              float32 `json:"mh_gspphot_upper" parquet:"name=mh_gspphot_upper, type=FLOAT"`
	DistanceGspphot             float32 `json:"distance_gspphot" parquet:"name=distance_gspphot, type=FLOAT"`
	DistanceGspphotLower        float32 `json:"distance_gspphot_lower" parquet:"name=distance_gspphot_lower, type=FLOAT"`
	DistanceGspphotUpper        float32 `json:"distance_gspphot_upper" parquet:"name=distance_gspphot_upper, type=FLOAT"`
	AzeroGspphot                float32 `json:"azero_gspphot" parquet:"name=azero_gspphot, type=FLOAT"`
	AzeroGspphotLower           float32 `json:"azero_gspphot_lower" parquet:"name=azero_gspphot_lower, type=FLOAT"`
	AzeroGspphotUpper           float32 `json:"azero_gspphot_upper" parquet:"name=azero_gspphot_upper, type=FLOAT"`
	AgGspphot                   float32 `json:"ag_gspphot" parquet:"name=ag_gspphot, type=FLOAT"`
	AgGspphotLower              float32 `json:"ag_gspphot_lower" parquet:"name=ag_gspphot_lower, type=FLOAT"`
	AgGspphotUpper              float32 `json:"ag_gspphot_upper" parquet:"name=ag_gspphot_upper, type=FLOAT"`
	EbpminrpGspphot             float32 `json:"ebpminrp_gspphot" parquet:"name=ebpminrp_gspphot, type=FLOAT"`
	EbpminrpGspphotLower        float32 `json:"ebpminrp_gspphot_lower" parquet:"name=ebpminrp_gspphot_lower, type=FLOAT"`
	EbpminrpGspphotUpper        float32 `json:"ebpminrp_gspphot_upper" parquet:"name=ebpminrp_gspphot_upper, type=FLOAT"`
	LibnameGspphot              string  `json:"libname_gspphot" parquet:"name=libname_gspphot, type=BYTE_ARRAY"`
}

func (schema GaiaInputSchema) GetId() string {
	return schema.Designation
}

func (schema GaiaInputSchema) GetCoordinates() (float64, float64) {
	return schema.RA, schema.Dec
}

func (schema GaiaInputSchema) FillMetadata(dst Metadata) {
	dst.(*Gaia).ID = schema.GetId()
	dst.(*Gaia).PhotGMeanFlux = sql.NullFloat64{Float64: schema.PhotGMeanFlux, Valid: true}
	dst.(*Gaia).PhotGMeanFluxError = sql.NullFloat64{Float64: float64(schema.PhotGMeanFluxError), Valid: true}
	dst.(*Gaia).PhotGMeanMag = sql.NullFloat64{Float64: float64(schema.PhotGMeanMag), Valid: true}
	dst.(*Gaia).PhotBpMeanFlux = sql.NullFloat64{Float64: float64(schema.PhotBpMeanFlux), Valid: true}
	dst.(*Gaia).PhotBpMeanFluxError = sql.NullFloat64{Float64: float64(schema.PhotBpMeanFluxError), Valid: true}
	dst.(*Gaia).PhotBpMeanMag = sql.NullFloat64{Float64: float64(schema.PhotBpMeanMag), Valid: true}
	dst.(*Gaia).PhotRpMeanFlux = sql.NullFloat64{Float64: float64(schema.PhotRpMeanFlux), Valid: true}
	dst.(*Gaia).PhotRpMeanFluxError = sql.NullFloat64{Float64: float64(schema.PhotRpMeanFluxError), Valid: true}
	dst.(*Gaia).PhotRpMeanMag = sql.NullFloat64{Float64: float64(schema.PhotRpMeanMag), Valid: true}
}

func (schema GaiaInputSchema) FillMastercat(dst *Mastercat, ipix int64) {
	dst.ID = schema.GetId()
	dst.Ra = schema.RA
	dst.Dec = schema.Dec
	dst.Cat = "gaia"
	dst.Ipix = ipix
}

func (gaia Gaia) GetId() string {
	return gaia.ID
}

func (gaia Gaia) GetCatalog() string {
	return "gaia"
}

func (m InsertGaiaParams) GetId() string {
	return m.ID
}

func (q *Queries) InsertGaiaWithoutParams(ctx context.Context, arg Gaia) error {
	_, err := q.db.ExecContext(
		ctx,
		insertGaia,
		arg.ID,
		arg.PhotGMeanFlux,
		arg.PhotGMeanFluxError,
		arg.PhotGMeanMag,
		arg.PhotBpMeanFlux,
		arg.PhotBpMeanFluxError,
		arg.PhotBpMeanMag,
		arg.PhotRpMeanFlux,
		arg.PhotRpMeanFluxError,
		arg.PhotRpMeanMag,
	)
	return err
}

func (m GetGaiaFromPixelsRow) GetId() string {
	return m.ID
}

func (m GetGaiaFromPixelsRow) GetCoordinates() (float64, float64) {
	return m.Ra, m.Dec
}

func (m GetGaiaFromPixelsRow) GetCatalog() string {
	return "gaia"
}
