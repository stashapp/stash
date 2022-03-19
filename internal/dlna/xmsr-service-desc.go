package dlna

// from https://github.com/rclone/rclone
// Copyright (C) 2012 by Nick Craig-Wood http://www.craig-wood.com/nick/

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

const xmsMediaReceiverServiceDescription = `<?xml version="1.0" ?>
<scpd xmlns="urn:schemas-upnp-org:service-1-0">
	<specVersion>
		<major>1</major>
		<minor>0</minor>
	</specVersion>
	<actionList>
		<action>
			<name>IsAuthorized</name>
			<argumentList>
				<argument>
					<name>DeviceID</name>
					<direction>in</direction>
					<relatedStateVariable>A_ARG_TYPE_DeviceID</relatedStateVariable>
				</argument>
				<argument>
					<name>Result</name>
					<direction>out</direction>
					<relatedStateVariable>A_ARG_TYPE_Result</relatedStateVariable>
				</argument>
			</argumentList>
		</action>
		<action>
			<name>RegisterDevice</name>
			<argumentList>
				<argument>
					<name>RegistrationReqMsg</name>
					<direction>in</direction>
					<relatedStateVariable>A_ARG_TYPE_RegistrationReqMsg</relatedStateVariable>
				</argument>
				<argument>
					<name>RegistrationRespMsg</name>
					<direction>out</direction>
					<relatedStateVariable>A_ARG_TYPE_RegistrationRespMsg</relatedStateVariable>
				</argument>
			</argumentList>
		</action>
		<action>
			<name>IsValidated</name>
			<argumentList>
				<argument>
					<name>DeviceID</name>
					<direction>in</direction>
					<relatedStateVariable>A_ARG_TYPE_DeviceID</relatedStateVariable>
				</argument>
				<argument>
					<name>Result</name>
					<direction>out</direction>
					<relatedStateVariable>A_ARG_TYPE_Result</relatedStateVariable>
				</argument>
			</argumentList>
		</action>
	</actionList>
	<serviceStateTable>
		<stateVariable sendEvents="no">
			<name>A_ARG_TYPE_DeviceID</name>
			<dataType>string</dataType>
		</stateVariable>
		<stateVariable sendEvents="no">
			<name>A_ARG_TYPE_Result</name>
			<dataType>int</dataType>
		</stateVariable>
		<stateVariable sendEvents="no">
			<name>A_ARG_TYPE_RegistrationReqMsg</name>
			<dataType>bin.base64</dataType>
		</stateVariable>
		<stateVariable sendEvents="no">
			<name>A_ARG_TYPE_RegistrationRespMsg</name>
			<dataType>bin.base64</dataType>
		</stateVariable>
		<stateVariable sendEvents="yes">
			<name>AuthorizationGrantedUpdateID</name>
			<dataType>ui4</dataType>
		</stateVariable>
		<stateVariable sendEvents="yes">
			<name>AuthorizationDeniedUpdateID</name>
			<dataType>ui4</dataType>
		</stateVariable>
		<stateVariable sendEvents="yes">
			<name>ValidationSucceededUpdateID</name>
			<dataType>ui4</dataType>
		</stateVariable>
		<stateVariable sendEvents="yes">
			<name>ValidationRevokedUpdateID</name>
			<dataType>ui4</dataType>
		</stateVariable>
	</serviceStateTable>
</scpd>`
