// validate IP address
function ValidateIPaddress(ipaddress) {  
  if ( /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(ipaddress) ) {  
    return true;
  }

  return false;
}  

// show the appropriate form on the page
function showTask(evt, cityName, fieldName) {
  var i, tabcontent, tablinks;
  tabcontent = document.getElementsByClassName("tabcontent");
  for (i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }
  tablinks = document.getElementsByClassName("tablinks");
  for (i = 0; i < tablinks.length; i++) {
    tablinks[i].className = tablinks[i].className.replace(" active", "");
  }
  document.getElementById(cityName).style.display = "block";
  document.getElementsByName(fieldName)[0].focus();
  evt.currentTarget.className += " active";
}

// puts a timestamp and a message in the results window
function printResult(message) {

        resultDiv = document.getElementById("results");

	dateTime = getNow();

	resultDiv.style.display = "block";
        document.getElementById("resultsControl").style.display = "block";

        if ( resultDiv.innerHTML == "" ) { 
                resultDiv.style.backgroundColor="light-grey";
                resultDiv.style.boxShadow="0 0 20px 1px rgba(0,0,0,.1)";
        }   

        resultDiv.innerHTML = "<div class='resultTime'>" + dateTime + "</div> " + message + "<br>" + resultDiv.innerHTML;

}

// padding single digit h:m:s with 0. :( why is this a thing?
function addZero(i) {
  if (i < 10) {
    i = "0" + i;
  }
  return i;
}

// just getting date and time;
function getNow() {

        var today = new Date();
        var date = today.getFullYear()+'-'+(today.getMonth()+1)+'-'+today.getDate();
        var time = addZero(today.getHours()) + ":" + addZero(today.getMinutes()) + ":" + addZero(today.getSeconds());
        var dateTime = date+' '+time;

	return dateTime;
}


// make the ajax call to the endpoint that the go binary is listening for
// /ping
// /port
// /dns
// /ssh
// /trace
// method will update an id responseText which will then get pulled into 

function makeAJAXCall(endpoint, data) {

	var http = new XMLHttpRequest();
	var url = endpoint;
	var params = data;
	http.open('POST', url, false)

	console.log(url + " " + params);
	
	//Send the proper header information along with the request
	http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
	
	http.onreadystatechange = function() {//Call a function when the state changes.
	    if(  http.status == 200) {
		console.log(url + " response:" + http.responseText);
		document.getElementById("responseText").innerHTML = http.responseText;
	    } 
	}
	http.send(params);
}

// function executed when Ping button is pressed
// todo: validate the input that we receive from the field. ensure it's IP or fqdn.
function goPing(host) {

	// trim whitespace off the ends snd get value
	hostValue = document.getElementsByName(host)[0].value.trim();
	
	// validate the input isn't blank
	if ( hostValue == "" ) {
		printResult("Please enter a value.");
		return;
	} 

	// validate IP
	if ( !ValidateIPaddress(hostValue) ) {
		printResult("Error: Invalid IP Address format. Format should be x.x.x.x");
		return false;
	}

	params = "pingHost=" + hostValue;

	makeAJAXCall('/ping', params);

	// log the value in the responseText element
	console.log(document.getElementById("responseText").innerHTML);
	pingResult = document.getElementById("responseText").innerHTML;

	printResult("PING " + pingResult);
}

// functon exxecuted when Port test button is pressed
// 
function goPortTest(host, port) {

	// take whitespace off the ends and get value
	hostValue = document.getElementsByName(host)[0].value.trim();
	portValue = document.getElementsByName(port)[0].value.trim();

	// check if input is blank
	if ( hostValue == "" || portValue == "" ) {
		printResult("Please provide a host and a port.");
		return;
	}

	// if there's a space
	if ( hostValue.indexOf(' ') >= 0 || portValue.indexOf(' ') >= 0 ) {
		printResult("Spaces are not allowed for host or port.");
		return;
	}

	if ( isNaN(portValue) ) {
		printResult("Port must be a number. '" + portValue + "'");
		return;
	}

	params = "Host=" + hostValue + "&Port=" + portValue;

	console.log("Params: " + params);

	makeAJAXCall('/port', params);

	console.log(document.getElementById("responseText").innerHTML);
	testResult = document.getElementById("responseText").innerHTML;

	printResult("PORT " + testResult);
}

function goDNSTest(fqdn) {

	// take whitespaces off the ends and get value
	fqdnValue = document.getElementsByName(fqdn)[0].value.trim();

	if ( fqdnValue == "" ) {
		printResult("Please provide an fqdn to lookup.");
		return;
	}

	if ( fqdnValue.indexOf(' ') >= 0 ) {
		printResult("Spaces are not allowed in the fqdn.");
		return;
	}

	params = "fqdn=" + fqdnValue

	console.log("Params: " + params);

	makeAJAXCall('/dns', params);

	console.log(document.getElementById("responseText").innerHTML);
	lookupResult = document.getElementById("responseText").innerHTML;

	printResult("LOOKUP " + fqdnValue + " - " +  lookupResult);

}
// function executed when ssh test button is pressed
function goSSHTest(sshHost, sshUser, sshPass, sshPrivKey) {

	// take whitespace off the ends and get value
	sshHostValue = document.getElementsByName(sshHost)[0].value.trim();
	sshUserValue = document.getElementsByName(sshUser)[0].value.trim();
	sshPassValue = document.getElementsByName(sshPass)[0].value.trim();
	sshPrivKeyValue = document.getElementsByName(sshPrivKey)[0].value.trim();

	// validate the input
	if ( sshHostValue == "" || sshUserValue == "" ) {
		printResult("Please provide a host and a username.");
		return;
	}

	if ( sshPassValue == "" && sshPrivKeyValue == "" ) {
		printResult("Password or private key must be provided. Private keys and passwords are not stored on disk.");
		return;
	}

	printResult(sshHostValue + " " + sshUserValue);


}

function goTraceHost(traceHost) {

	// take whitespace off the ends and get value
	traceHostValue = document.getElementsByName(traceHost)[0].value.trim();

	if ( traceHostValue == "" ) {
		printResult("Please provide a host to trace.");
		return;
	}

	if ( traceHostValue.indexOf(' ') >= 0 ) {
		printResult("Spaces are not allowed for host.");
		return;
	}

	printResult(traceHostValue);

}
	

