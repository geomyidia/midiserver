(defmodule sp-midi
  (export all))

(include-lib "sp_midi.hrl")

(defun lib-file () (LIB))

(defun nif-loaded? () (sp_midi:is_nif_loaded))
(defun nif-initialised? () (sp_midi:is_nif_initialized))

(defun initialise () (sp_midi:midi_init))
(defun deinitialise () (sp_midi:midi_deinit))

(defun send (a b) (sp_midi:midi_send a b))
(defun flush () (sp_midi:midi_flush))

(defun inputs () (sp_midi:midi_ins))
(defun outputs () (sp_midi:midi_outs))
(defun refresh () (sp_midi:midi_refresh_devices))

(defun have-pid? () (sp_midi:have_my_pid))

(defun current-ms () (sp_midi:get_current_time_microseconds))
(defun set-log-level (a) (sp_midi:set_log_level a))
(defun set-pid (a) (sp_midi:set_this_pid a))
