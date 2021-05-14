
async function displayPreferences(matches, div) {

    let preferences = await getUrl('api/preference')

    if (validatePreferences(preferences)) {
        preferences = createDefaultPreferences()
    }

    div.innerHTML = `
        <div class="container h-100">
            <div class="row px-3" id="alert"></div>
            <div class="text-light text-center h3 mb-2">Preferences</div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Width</label>
                <div class="col-sm-10" id="dimensionWidthId">
                    ${createNumberInput(preferences.dimensions.width)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Length</label>
                <div class="col-sm-10" id="dimensionLengthId">
                    ${createNumberInput(preferences.dimensions.length)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Units</label>
                <div class="col-sm-10">
                    <select class="form-control" id="dimensionUnitsId" >
                        <option value="in" ${preferences.dimensions.units === 'in' ? 'selected' : ''}>Inches</option>
                        <option value="cm" ${preferences.dimensions.units === 'cm' ? 'selected' : ''}>Centimeters</option>
                    </select>
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Width Axis</label>
                <div class="col-sm-10">
                    <select class="form-control" id="orientationWidthId" >
                        <option id="orientationWidthIdX"  value="x" ${preferences.orientation.width === 'x' ? 'selected' : ''}>X</option>
                        <option id="orientationWidthIdY"  value="Y" ${preferences.orientation.width === 'y' ? 'selected' : ''}>Y</option>
                        <option id="orientationWidthIdZ"  value="Z" ${preferences.orientation.width === 'z' ? 'selected' : ''}>Z</option>
                    </select>
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Length Axis</label>
                <div class="col-sm-10">
                    <select class="form-control" id="orientationLengthId" >
                        <option id="orientationLengthIdX"  value="x" ${preferences.orientation.length === 'x' ? 'selected' : ''}>X</option>
                        <option id="orientationLengthIdY"  value="Y" ${preferences.orientation.length === 'y' ? 'selected' : ''}>Y</option>
                        <option id="orientationLengthIdZ"  value="Z" ${preferences.orientation.length === 'z' ? 'selected' : ''}>Z</option>
                    </select>
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Tolerance</label>
                <div class="col-sm-10" id="toleranceId">
                    ${createNumberInput(preferences.tolerance)}
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Update Rate</label>
                <div class="col-sm-10" id="updateRateId">
                    ${createNumberInput(preferences.updateRate)}
                </div>
            </div>
            <div class="custom-control custom-switch">
              <input type="checkbox" class="custom-control-input" id="orientationInvertRollId" ${preferences.orientation.invertRoll ? 'checked' : ''}>
              <label class="custom-control-label text-light h5" for="orientationInvertRollId">Invert Roll</label>
            </div>
            <div class="custom-control custom-switch">
              <input type="checkbox" class="custom-control-input" id="orientationInvertPitchId" ${preferences.orientation.invertPitch ? 'checked' : ''}>
              <label class="custom-control-label text-light h5" for="orientationInvertPitchId">Invert Pitch</label>
            </div>
            <button type="button" class="btn btn-lg btn-primary btn-block mt-5" onclick="savePreferences()">Save</button>
            <!-- TODO Add button to export the database -->
        </div>
    `;
}

function createDefaultPreferences() {
    return {
        dimensions: {
            width: 96.0,
            length: 240.0,
            units: "in",
        },
        orientation: {
            length: "y",
            width: "x",
            invertPitch: false,
            invertRoll: false
        },
        tolerance: 0.1,
        updateRate: 1.0
    }
}

async function savePreferences() {

    const preferences = createDefaultPreferences()
    preferences.dimensions.width = parseFloat(document.getElementById('dimensionWidthId').firstElementChild.value)
    preferences.dimensions.length = parseFloat(document.getElementById('dimensionLengthId').firstElementChild.value)
    preferences.dimensions.units = document.getElementById('dimensionUnitsId').value
    preferences.tolerance = parseFloat(document.getElementById('toleranceId').firstElementChild.value)
    preferences.updateRate = parseFloat(document.getElementById('updateRateId').firstElementChild.value)

    preferences.orientation.width = document.getElementById('orientationWidthIdX').selected ? 'x' : preferences.orientation.width
    preferences.orientation.width = document.getElementById('orientationWidthIdY').selected ? 'y' : preferences.orientation.width
    preferences.orientation.width = document.getElementById('orientationWidthIdZ').selected ? 'z' : preferences.orientation.width

    preferences.orientation.length = document.getElementById('orientationLengthIdX').selected ? 'x' : preferences.orientation.length
    preferences.orientation.length = document.getElementById('orientationLengthIdY').selected ? 'y' : preferences.orientation.length
    preferences.orientation.length = document.getElementById('orientationLengthIdZ').selected ? 'z' : preferences.orientation.length

    preferences.orientation.invertRoll = document.getElementById('orientationInvertRollId').checked
    preferences.orientation.invertPitch = document.getElementById('orientationInvertPitchId').checked

    console.log(preferences)

    const validationResult = validatePreferences(preferences)
    if (validationResult) {
        displayAlert('alert-danger', validationResult)
        return
    }

    await putUrl('api/preference', preferences)

    displayAlert('alert-success', 'Update success', true)
}

function validatePreferences(preferences) {

    if (!preferences.dimensions.width || (preferences.dimensions.width <= 0.0)) {
        return 'Width must be greater than zero'
    }

    if (!preferences.dimensions.length || (preferences.dimensions.length <= 0.0)) {
        return 'Length must be greater than zero'
    }

    if (!preferences.tolerance || (preferences.tolerance <= 0.0)) {
        return 'Tolerance must be greater than zero'
    }

    if (!preferences.updateRate || (preferences.updateRate <= 0.0)) {
        return 'Update Rate must be greater than zero'
    }

    if (!preferences.orientation.width) {
        return "Width Axis is not set"
    }

    if (!preferences.orientation.length) {
        return "Length Axis is not set"
    }

    if (preferences.orientation.width === preferences.orientation.length) {
        return "Width and Length Axis must not be set to the same value"
    }

    return null
}
