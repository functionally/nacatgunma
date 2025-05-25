
'use strict'

import * as DagCbor from "@ipld/dag-cbor"
import * as Cbor  from "cbor2"

import { Buffer } from "buffer"
import { CID } from "multiformats/cid"
import { DataSet} from "vis-data"
import { Network} from "vis-network"


const textDecoder = new TextDecoder()

function reportError(message) {
  window.status = message
  console.error(message)
}

function ensureTrailingSlash(url) {
  return url.endsWith("/") ? url : str + "/"
}


async function fetchScriptUtxos(followup) {
  if (uiBlockfrostToken.value == null || uiBlockfrostToken.value == "") {
    alert("A Blockfrost token is required.")
    return
  } else if (uiBlockfrostUrl.value == null || uiBlockfrostUrl.value == "") {
    alert("The Blockfrost URL is required.")
    return
  } else if (uiScriptAddress.value == null | uiScriptAddress.value == "") {
    alert("The Nacatgunma script address is required.")
    return
  }
  const xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        const res = JSON.parse(this.responseText)
        followup(res)
      } else {
        reportError("Blockfrost status: " + this.status)
      }
    }
  }
  xhttp.open("GET", ensureTrailingSlash(uiBlockfrostUrl.value) + "addresses/" + uiScriptAddress.value + "/utxos")
  xhttp.setRequestHeader("project_id", uiBlockfrostToken.value)
  xhttp.setRequestHeader("Accept", "application/json")
  xhttp.send()
}

async function fetchIpldCbor(cid, followup) {
  if (uiIpfsGateway.value == null || uiIpfsGateway.value =="") {
    alert("An IPFS gateway URL is required.")
    return
  }
  const xhttp = new XMLHttpRequest()
  xhttp.responseType = "arraybuffer"
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4) {
      if (this.status == 200) {
        const contentType = this.getResponseHeader("Content-type")
        const response = this.response
        if (contentType == "application/vnd.ipld.dag-json") {
          const res = JSON.parse(textDecoder.decode(response))
          followup(res)
        } else if (contentType == "application/vnd.ipld.dag-cbor") {
          const res = DagCbor.decode(new Uint8Array(response))
          followup(res)
        } else {
          reportError("Failed to decode block header (Content-type = " + contentType + ").")
        }
      } else {
        reportError("IPFS gateway status: " + this.status)
      }
    }
  }
  xhttp.open("GET", ensureTrailingSlash(uiIpfsGateway.value) + "ipfs/" + cid)
  xhttp.setRequestHeader("Accept", "application/vnd.ipld.dag-cbor")
  xhttp.send()
}


function shortenLabel(label) {
  return label.slice(0, 5) + ".." + label.slice(-5)
}

function extractCid(x) {
  const xString = x.toString()
  try {
    CID.parse(xString)
    return xString
  } catch(err) {
  }
  try {
    return x["/"]
  } catch(err) {
  }
  return null
}

function createHTMLTitle(html) {
  var element = document.createElement("div")
  element.innerHTML = html
  return element
}

function utxoNode(utxo, level = 0) {
  const utxoId = utxo.tx_hash + "#" + utxo.tx_index 
  if (data.nodes.getIds().filter(id => id == utxoId).length > 0)
    return utxoId
  data.nodes.add({
    id: utxoId,
    title: createHTMLTitle("<table><tr><th>UTxO</th><td><code>" + utxoId + "</code></td></tr>"),
    label: shortenLabel(utxoId),
    color: "coral",
    shape: "ellipse",
//  level: level - 1,
  })
  return utxoId
}

function headerNode(headerId, tooltip, level = 0) {
  if (data.nodes.getIds().filter(id => id == headerId).length > 0)
    return headerId
  data.nodes.add({
    id: headerId,
    title: createHTMLTitle("<table><tr><th>Block header</th><td><code>" + headerId + "</code></td></tr>" + tooltip + "</table>"),
    label: shortenLabel(headerId),
    color: "cornflowerblue",
    shape: "box",
//  level,
  })
  return headerId
}

function bodyNode(bodyCid, tooltip, level = 0) {
  const bodyId = extractCid(bodyCid)
  if (bodyId == null)
    return null
  if (data.nodes.getIds().filter(id => id == bodyId).length > 0)
    return bodyId
  data.nodes.add({
    id: bodyId,
    title: createHTMLTitle("<table><tr><th>Block body</th><td><code>" + bodyId + "</code></td></tr>" + tooltip + "</table>"),
    label: shortenLabel(bodyId),
    color: "darkseagreen",
    shape: "box",
//  level,
    shapeProperties: {
      borderRadius: 0,
    },
  })
  return bodyId
}

function utxoEdge(utxoId, headerId) {
  const edgeId = utxoId + "|" + headerId
  if (data.edges.getIds().filter(id => id == edgeId).length > 0)
    return edgeId
  data.edges.add({
    id: edgeId,
    from: utxoId,
    to: headerId,
    color: "coral",
  })
  return edgeId
}

function acceptEdge(parentId, headerId) {
  const edgeId = parentId + "|" + headerId
  if (data.edges.getIds().filter(id => id == edgeId).length > 0)
    return edgeId
  data.edges.add({
    id: edgeId,
    from: parentId,
    to: headerId,
    color: "cornflowerblue",
  })
  return edgeId
}

function rejectEdge(parentId, headerId) {
  const edgeId = parentId + "|" + headerId
  if (data.edges.getIds().filter(id => id == edgeId).length > 0)
    return edgeId
  data.edges.add({
    id: edgeId,
    from: parentId,
    to: headerId,
    color: "crimson",
    dashes: true,
  })
  return edgeId
}

function bodyEdge(headerId, bodyId) {
  const edgeId = headerId + "|" + bodyId
  if (data.edges.getIds().filter(id => id == edgeId).length > 0)
    return edgeId
  data.edges.add({
    id: edgeId,
    from: headerId,
    to: bodyId,
    color: "darkseagreen",
  })
  return edgeId
}

function addBlock(headerCid, level = 0) {
  const headerId = extractCid(headerCid)
  if (headerId == null)
    return null
  fetchIpldCbor(headerId, function(result) {
    headerNode(
      headerId,
      "<tr><th>Issuer</th><td><code>" + result.Issuer + "</code></td></tr><tr><th>Comment</th><td><code>" + result.Payload.Comment + "</code></td></tr>",
      level
    )
    const bodyId = bodyNode(
      result.Payload.Body,
      "<tr><th>Schema</th><td><code>" + result.Payload.Schema + "</code></td></tr><tr><th>Media type</th><td><code>" + result.Payload.MediaType + "</code></td></tr>",
      level
    )
    bodyEdge(headerId, bodyId)
    if (level >= uiLevelLimit.value)
      return headerId
    result.Payload.Accept.forEach(function(accept) {
      const blockId = addBlock(accept, level + 1)
      acceptEdge(headerId, blockId)
    })
    result.Payload.Reject.forEach(function(reject) {
      const blockId = addBlock(reject, level + 1)
      reject(headerId, blockId)
    })
  })
  return headerId
}

function addUtxo(utxo) {
  const utxoId = utxoNode(utxo)
  const headerId = addBlock(utxo.headerString)
  utxoEdge(utxoId, headerId)
}

async function fetchTips() {
  fetchScriptUtxos(function(candidates) {
    const utxos = candidates.filter(function(utxo) {
      let found = uiFilterToken.value == ""
      const amount = utxo.amount
      for (let i in amount)
        if (amount[i].unit == uiFilterToken.value)
          found = true
      return found
    })
    const tips = utxos.map(function(utxo) {
      const cid = CID.decode(Cbor.decode(Buffer.from(utxo.inline_datum, "hex"))[1])
      return {
        tx_hash: utxo.tx_hash,
        tx_index: utxo.tx_index,
        headerCid: cid,
        headerString: cid.toString(),
      }
    })
    tips.forEach(addUtxo)
  })
}


export const data = {
  nodes: new DataSet(),
  edges: new DataSet(),
}

export function drawBlocks() {
  const options = {
    layout: {
      hierarchical: {
        direction: "RL",
        sortMethod: "directed",
      },
      improvedLayout: false,
    },
    physics: {
      enabled: true,
    },
    nodes: {
      font: {face: "monospace"},
    },
    edges: {
      arrows: {
        to: { enabled: true, scaleFactor: 1 },
      },
      smooth: true
    },
    interaction: { hover: true },
  }
  const network = new Network(uiCanvas, data, options)
  network.on("click", function (params) {
    if (params.nodes.length > 0) {
      const nodeId = params.nodes[0]
      try {
        CID.parse(nodeId)
        window.open(ensureTrailingSlash(uiIpldExplorer.value) + nodeId, "nacatgunma")
      } catch (err) {
        if (nodeId.includes("#"))
          window.open(ensureTrailingSlash(uiCardanoExplorer.value) + "transaction/" + nodeId.split('#')[0], "nacatgunma")
      }
    }
  })
  setTimeout(() => {
    network.setOptions({ physics: false })
  }, 15000)
}


const KEY_SCRIPT_ADDRESS = "scriptAddress"
const KEY_FILTER_TOKEN = "filterToken"
const KEY_BLOCKFROST_URL = "blockfrostUrl"
const KEY_BLOCKFROST_TOKEN = "blockfrostToken"
const KEY_IPFS_GATEWAY = "ipfsGateway"
const KEY_IPLD_EXPLORER = "ipldExplorer"
const KEY_CARDANO_EXPLORER = "cardanoExplorer"
const KEY_LEVEL_LIMIT = "levelLimit"

function setupPersistence(key, element, defaultValue, followup) {
  const value = localStorage.getItem(key)
  if (value) {
    element.value = value
  } else if (defaultValue) {
    element.value = defaultValue
  }
  element.addEventListener("change", function() {
    localStorage.setItem(key, element.value)
    if (followup)
      followup()
  })
}

function reset() {
  data.nodes =new DataSet()
  data.edges =new DataSet()
  fetchTips()
}


export async function initialize() {

  setupPersistence(KEY_SCRIPT_ADDRESS, uiScriptAddress, "addr1w8lyu0uj30gyytukg25ynfypvqlw7tt4duuu7lqd09qrnugm34xp8", reset)
  setupPersistence(KEY_FILTER_TOKEN, uiFilterToken, "30135f08305143796de4276083cc54e47fbcafb176df6b58ab3094464e6163617467756e6d61", reset)
  setupPersistence(KEY_BLOCKFROST_URL, uiBlockfrostUrl, "https://cardano-mainnet.blockfrost.io/api/v0/", reset)
  setupPersistence(KEY_BLOCKFROST_TOKEN, uiBlockfrostToken, null, reset)
  setupPersistence(KEY_IPFS_GATEWAY, uiIpfsGateway, "https://ipfs.io/", reset)
  setupPersistence(KEY_IPLD_EXPLORER, uiIpldExplorer, "https://explore.ipld.io/#/explore/")
  setupPersistence(KEY_CARDANO_EXPLORER, uiCardanoExplorer, "https://cardanoscan.io/")
  setupPersistence(KEY_LEVEL_LIMIT, uiLevelLimit, 100)

  drawBlocks()

  fetchTips()

}
