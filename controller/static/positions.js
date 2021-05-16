
function displayHome(matches, div) {
    fetch('api/position')
        .then(res => res.json())
        .then(positions => {
            const favorites = positions.filter(value => {
                return value.favorite;
            })
            div.innerHTML = `
                <div class="container h-100">
                    <a class="btn btn-lg btn-primary btn-block mb-4" role="button" href="#position/level"><div class="display-4">Level</div></a>
                    ${favorites.length === 0 ? `<div class="text-secondary text-center">You have no favorites</div>` : createPositionButtons(favorites)}
                    <a class="btn btn-lg btn-primary btn-block mt-4" role="button" href="#position">View All Positions</a>
                    <a class="btn btn-lg btn-success btn-block" role="button" href="#position/new">Save Current Position</a>
                </div>
            `;
        })
        .catch(err => console.error(err));
}

function displayPositionList(matches, div) {
    fetch('api/position')
        .then(res => res.json())
        .then(positions => {
            div.innerHTML = `
                <div class="container h-100">
                    <div class="text-light text-center h3">Positions</div>
                    ${positions.length === 0 ? `<div class="text-secondary text-center">You have no saved positions</div>` : createPositionButtons(positions)}
                </div>
            `;
        })
        .catch(err => console.error(err));
}

function createPositionButtons(positions) {
    return positions.map(item => `<a class="btn btn-lg btn-secondary btn-block" role="button" href="#position/${item.id}/recall">${item.name}</a>`.trim()).join('')
}
