
async function getUrls(urls) {
  return await Promise.all(
    urls.map(url => fetch(url)
      .then(r => r.json())
      .catch(err => console.error(err)))
  )
}

async function getUrl(url) {
  const results = await getUrls([url])
  return results[0]
}

function putUrl(url, body) {
  return fetch(url, {
    method: 'PUT',
    body: JSON.stringify(body),
  })
    .then(r => r.json())
    .then(r => {
      displayAlert('alert-success', 'Update success', true)
      return r
    })
    .catch(err => {
      console.error(err)
      displayAlert('alert-danger', 'Update failed', false)
    });
}

function postUrl(url, body) {
  console.log(`Creating ${url} with ${JSON.stringify(body)}`)
  return fetch(url, {
    method: 'POST',
    body: JSON.stringify(body),
  })
    .then(r => r.json())
    .then(r => {
      displayAlert('alert-success', 'Create success', true)
      return r
    })
    .catch(err => {
      console.error(err)
      displayAlert('alert-danger', 'Create failed', false)
    });
}

function deleteUrl(url) {
  console.log(`Deleting ${url}`)
  return fetch(url, {
    method: 'DELETE'
  })
    .then(r => {
      displayAlert('alert-success', 'Delete success', true)
      return r
    })
    .catch(err => {
      console.error(err)
      displayAlert('alert-danger', 'Delete failed', false)
    });
}

function displayAlert(type, text, autoDismiss) {
  let alertDiv = document.getElementById('alert');
  alertDiv.innerHTML = `
    <div class="alert ${type} alert-dismissible fade show w-100" role="alert">
      ${text}
      <button type="button" class="close" data-dismiss="alert" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
    </div>
  `;
  if (autoDismiss) {
    setTimeout(function () {
      alertDiv.innerHTML = '';
    }, 2000);
  }
}

function median(values) {
  if (values.length === 0) {
    return 0;
  }
  values.sort(function (a, b) {
    return a - b;
  });
  const half = Math.floor(values.length / 2);
  if (values.length % 2) {
    return values[half];
  }
  return (values[half - 1] + values[half]) / 2.0;
}
