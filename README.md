a bit of interview prep. I asked meta AI to give me a protocol parsing problem.

```
Here is a protocol coding challenge:
Protocol Definition
The "Simple Messaging Protocol" (SMP) is a binary protocol used for sending messages between devices. It has the following format:
Field Name
Type
Length (bytes)
Description
Message ID
uint16
2
Unique ID of the message
Message Type
uint8
1
Type of message (1=data, 2=control)
Payload Length
uint16
2
Length of the payload in bytes
Payload
bytes
variable
The actual message payload
CRC
uint16
2
Cyclic Redundancy Check of the payload
Challenge
Write a function in your preferred programming language (e.g. Python, C++, Java) that takes a byte slice as input and parses it according to the SMP protocol. The function should return a struct or object containing the parsed fields.
Example Input
The following byte slice represents a valid SMP message:
[0x12, 0x34, 0x01, 0x00, 0x05, 0x00, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x1a, 0x2b]
Expected Output
The parsed struct or object should contain the following fields:
{
  "Message ID": 0x1234,
  "Message Type": 0x01,
  "Payload Length": 0x0005,
  "Payload": [0x68, 0x65, 0x6c, 0x6c, 0x6f],
  "CRC": 0x1a2b
}
Good Luck!
Let me know if you need any clarification or help with the challenge.
```