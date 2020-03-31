// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/***************************************************************************
 *  Copyright 1991, 1992, 1993, 1994, 1995, 1996, 2001, 2002               *
 *    David R. Hill, Leonard Manzara, Craig Schock                         *
 *                                                                         *
 *  This program is free software: you can redistribute it and/or modify   *
 *  it under the terms of the GNU General Public License as published by   *
 *  the Free Software Foundation, either version 3 of the License, or      *
 *  (at your option) any later version.                                    *
 *                                                                         *
 *  This program is distributed in the hope that it will be useful,        *
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of         *
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the          *
 *  GNU General Public License for more details.                           *
 *                                                                         *
 *  You should have received a copy of the GNU General Public License      *
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.  *
 ***************************************************************************/
// 2014-09
// This file was copied from Gnuspeech and modified by Marcelo Y. Matuda.

// 2019-02
// This is a port to golang of the C++ Gnuspeech port by Marcelo Y. Matuda

package v2

import (
	"strconv"

	"github.com/goki/ki/bitflag"
)

const modelConfigFn = "/trm_control_model.config"
const resonanceConfigFn = "/trm.config"
const voiceFilePrefix = "/voice_"

type Control struct {
	Model       *Model
	Sequence    Sequence
	ModelConfig ModelConfig
	TrmConfig   TrmConfig
}

func (ctrl *Control) Init(path string, model *Model) {
	ctrl.Model = model
	ctrl.Model = model
	ctrl.Sequence.Init(path, model)
	ctrl.LoadConfigs("")
}

//
func (ctrl *Control) LoadConfigs(path string) {

	ctrl.ModelConfig.Load(path + modelConfigFn)

	resonanceConfigPath := path + resonanceConfigFn

	voiceConfigPath := path + voiceFilePrefix + ".config"

	ctrl.TrmConfig.Load(resonanceConfigPath, voiceConfigPath)
}

func (ctrl *Control) InitUtterance() {
	rc := NewTrmConfig()
	mc := ctrl.ModelConfig
	if rc.OutputRate != 22050.0 && rc.OutputRate != 44100.0 {
		rc.OutputRate = 44100.0
	}
	if rc.VtlOffset+rc.VocalTractLength < 15.9 {
		rc.OutputRate = 44100.0
	}

	ctrl.Sequence.PitchMean = ctrl.ModelConfig.PitchOffset + rc.RefGlottalPitch
	ctrl.Sequence.GlobalTempo = mc.Tempo
	ctrl.SetIntonation(mc.Intonation)
	ctrl.Sequence.Drift.SetUp(mc.DriftDeviation, mc.ControlRate, mc.DriftLowCutoff)
	ctrl.Sequence.SetRadiusCoefs(rc.RadiusCoefs)

	trmParamStream<<
		rc.outputRate<<'\n'<<
		mc.ControlRate<<'\n'<<
		rc.volume<<'\n'<<
		rc.channels<<'\n'<<
		rc.balance<<'\n'<<
		rc.waveform<<'\n'<<
		rc.glottalPulseTp<<'\n'<<
		rc.glottalPulseTnMin<<'\n'<<
		rc.glottalPulseTnMax<<'\n'<<
		rc.breathiness<<'\n'<<
		rc.vtlOffset + rc.TractLength<<'\n'<< // tube length
		rc.temperature<<'\n'<<
		rc.lossFactor<<'\n'<<
		rc.apertureRadius<<'\n'<<
		rc.mouthCoef<<'\n'<<
		rc.noseCoef<<'\n'<<
		rc.noseRadius[1]<<'\n'<<
		rc.noseRadius[2]<<'\n'<<
		rc.noseRadius[3]<<'\n'<<
		rc.noseRadius[4]<<'\n'<<
		rc.noseRadius[5]<<'\n'<<
		rc.throatCutoff<<'\n'<<
		rc.throatVol<<'\n'<<
		rc.modulation<<'\n'<<
		rc.mixOffset<<'\n'
}

// Chunks are separated by /c.
// There is always one /c at the begin and another at the end of the string.
func (ctrl *Control) CalcChunks(text string) int {
	tmp := 0
	idx := 0
	for text[idx] != '0' {
		if (text[idx] == '/') && (text[idx+1] == 'c') {
			tmp++
			idx += 2
		} else {
			idx++
		}
	}
	tmp--
	if tmp < 0 {
		tmp = 0
	}
	return tmp
}

// NextChunk returns the position of the next /c (the position of the /).
func (ctrl *Control) NextChunk(text string) int {
	idx := 0
	for text[idx] != '0' {
		if (text[idx] == '/') && (text[idx+1] == 'c') {
			return idx
		} else {
			idx++
		}
	}
	return 0
}

// ValidPosture
func (ctrl *Control) ValidPosture(token string) bool {
	i, err := strconv.Atoi(string(token[0]))
	if err != nil {
		return false
	}

	if i >= 0 && i <= 9 {
		return true
	} else {
		return ctrl.Model.Postures.PostureTry(token) != nil
	}
}

// SetIntonation
func (ctrl *Control) SetIntonation(intonation int64) {
	ctrl.Sequence.MicroFlag = 0
	ctrl.Sequence.MacroFlag = 0
	ctrl.Sequence.SmoothInton = 0 // Macro and not smooth is not working.
	ctrl.Sequence.DriftFlag = 0
	ctrl.Sequence.TgUseRandom = false

	if bitflag.Has(intonation, int(IntonationMicro)) {
		ctrl.Sequence.MicroFlag = 1
	}

	if bitflag.Has(intonation, int(IntonationMacro)) {
		ctrl.Sequence.MacroFlag = 1
		ctrl.Sequence.SmoothInton = 1 // Macro and not smooth is not working.
	}

	// Macro and not smooth is not working.
	// if bitflag.Has(intonation, int(IntonationSmooth)) {
	// 	ctrl.Sequence.SetSmoothIntonation(1)
	// }

	if bitflag.Has(intonation, int(IntonationDrift)) {
		ctrl.Sequence.DriftFlag = 1
	}

	if bitflag.Has(intonation, int(IntonationRandom)) {
		ctrl.Sequence.TgUseRandom = true
	}
}
