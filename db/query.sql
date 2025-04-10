SELECT DISTINCT
tps.study_id,
tps.submitted_study_name as StudyName,
tps.faculty_sponsor,
tps.programme,
tps.sanger_sample_id,
tps.supplier_name,
tps.manifest_created,
tps.manifest_uploaded,
tps.labware_received,
sr.labware_human_barcode as 'Plate/Tube',
tps.order_made,
tps.library_start,
tps.library_complete,
DATEDIFF(tps.library_complete,tps.library_start) as LibraryTime,
COALESCE(pm.id_run,
pbr.id_pac_bio_run_lims,
ofc.experiment_name) as RunID,
COALESCE(rlm.instrument_model, pbw.instrument_type, ofc.instrument_name) as Platform,
COALESCE(fc.pipeline_id_lims,
pbr.pipeline_id_lims,
ofc.pipeline_id_lims) as Pipeline,
tps.sequencing_run_start,
tps.sequencing_qc_complete,
DATEDIFF(tps.sequencing_qc_complete,tps.sequencing_run_start) as SequencingTime,
COALESCE(fc.manual_qc, pbm.qc) as qcPass

FROM mlwh_reporting.seq_ops_tracking_per_sample tps
JOIN sample sa on sa.id_sample_lims = tps.id_sample_lims
JOIN stock_resource sr on sr.id_sample_tmp = sa .id_sample_tmp

LEFT JOIN iseq_flowcell fc on fc.id_sample_tmp = sa.id_sample_tmp
LEFT JOIN iseq_product_metrics pm on pm.id_iseq_flowcell_tmp = fc.id_iseq_flowcell_tmp
LEFT JOIN iseq_run_lane_metrics rlm on rlm.id_run = pm.id_run

LEFT JOIN pac_bio_run pbr on pbr.id_sample_tmp = sa.id_sample_tmp
LEFT JOIN pac_bio_product_metrics pbm on pbm.id_pac_bio_tmp = pbr.id_pac_bio_tmp
LEFT JOIN pac_bio_run_well_metrics pbw on pbw.id_pac_bio_product = pbm.id_pac_bio_product

LEFT JOIN oseq_flowcell ofc on ofc.id_sample_tmp = sa.id_sample_tmp

WHERE tps.manifest_created >= DATE_SUB(NOW(), INTERVAL 2 YEAR)
AND tps.faculty_sponsor IN ('Carl Anderson' 'Emma Davenport','Matthew Hurles','Hilary Martin','Gosia Trynka', 'Ben Lehner', 'Leopold Parts', 'Jussi Taipale')

ORDER BY COALESCE(tps.manifest_created, tps.manifest_uploaded), tps.study_id, tps.sanger_sample_id, RunID; 