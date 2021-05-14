
async function displayCalibration(matches, div) {

    const calibration = await getUrl('api/calibration')

    div.innerHTML = `
        <div class="container h-100">
            <div class="row px-3" id="alert"></div>
            <div class="text-light text-center h3 mb-2">Calibration</div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Roll</label>
                <div class="col-sm-10" id="rollId">
                    ${createNumberInput(calibration.roll)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Pitch</label>
                <div class="col-sm-10" id="pitchId">
                    ${createNumberInput(calibration.pitch)}
                </div>
            </div>
            <button type="button" class="btn btn-lg btn-primary btn-block mt-5" onclick="refreshCalibration()">Refresh</button>
            <button type="button" class="btn btn-lg btn-danger btn-block" data-toggle="modal" data-target="#saveModal">Save</button>
            <div class="modal fade" id="saveModal" tabindex="-1" role="dialog" aria-labelledby="saveModalLabel" aria-hidden="true">
              <div class="modal-dialog" role="document">
                <div class="modal-content">
                  <div class="modal-header">
                    <h5 class="modal-title" id="saveModalLabel">Update Calibration</h5>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                      <span aria-hidden="true">&times;</span>
                    </button>
                  </div>
                  <div class="modal-body">Are you sure you want to update the calibration?</div>
                  <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                    <button type="button" class="btn btn-danger" onclick="saveCalibration()" data-dismiss="modal">Save</button>
                  </div>
                </div>
              </div>
            </div>
        </div>
    `;
}

async function refreshCalibration() {
    const uncorrected = await getUrl('api/uncorrected')
    document.getElementById('rollId').innerHTML = createNumberInput(uncorrected.roll)
    document.getElementById('pitchId').innerHTML = createNumberInput(uncorrected.pitch)
}

async function saveCalibration() {
    const roll = parseFloat(document.getElementById('rollId').firstElementChild.value)
    const pitch = parseFloat(document.getElementById('pitchId').firstElementChild.value)
    await putUrl('api/calibration', { roll: roll, pitch: pitch })
    displayAlert('alert-success', 'Update success', true)
}
