# MOH Apps User Management


This is the User Management Application for the National AIDS Commision(NAC) Surveillance System.

## Roles
A user of the system can have one of the following roles:
1. Adherence Counselor
2. Peer Navigator
3. SR

Users can belong to one of the following organizations:
1. CSO
2. BFLA
3. NAC
4. MOHW

Users at NAC and MOHW will be able to access data for both CSO and BFLA in summary form. They will
not have permission to add or modify data.

Users at CSO and BFLA will only be able to add or modify data for their respective organization.
They will be able to view data entered at any organization. 

## Adherence Counselor
An Adherence Counselor will only be able to edit information in an Adherence Form. Data entered by the
Peer Navigator will be collated into an Adherence Form for these counselors. They should not be allowed
to edit Peer Navigator Forms.

## Peer Navigator
A Peer Navigator can only edit data in the Peer Navigator Intake Form. They can not edit Adherence 
Intake Forms. However, they can view the data entered by the Adherence Counselors.

## SR
An SR can view all the data entered in the system, but can not edit data.



## User Object Structure
Internally, a user object has the following form:

```
{
id: string;
firstName: string;
lastName: string;
email: string;
org: CSO | BFLA | NAC | NOHW;
role: Adherence_Counselor | Peer_Navigator | SR 
enabled: boolean;
createdBy: string; // The email of the user who created this record
updatedBy: string; // The email of the user who last updated this record
updatedDate: Date;
createdDate: Date;
}

```
