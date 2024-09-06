(defmodule midiserver-app
  (behaviour application)
  ;; app implementation
  (export
   (start 2)
   (stop 1)))

;;; --------------------------
;;; application implementation
;;; --------------------------

(defun start (_type _args)
  (logger:set_application_level 'midiserver 'all)
  (logger:info "Starting midiserver application ...")
  (midiserver-sup:start_link))

(defun stop (_state)
  (midiserver-sup:stop)
  'ok)
