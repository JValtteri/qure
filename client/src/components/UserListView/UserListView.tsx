import { useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';

import Frame from '../common/Frame/Frame';
import GenericTable from '../common/GenericTable/GenericTable';
import { adminDeleteUser, listAllClients, type ClientResponse } from '../../api/api';
import UserInspectCard from '../UserInspectCard/UserInspectCard';
import ConfirmDeleteDialog from '../ConfirmDeleteDialog/ConfirmDeleteDialog';


const loadingClientList = signal(false);

interface Props {
    active: boolean;
    setShowPopup: (value: React.SetStateAction<boolean>) => void
    setPopupMessage: (value: React.SetStateAction<string>) => void;
}

function UserListView({active, setShowPopup, setPopupMessage}: Props) {
    const [data, setData] = useState(new Array<ClientResponse>());
    const [targetClient, setTargetClient] = useState({} as ClientResponse); // ClientResponse
    const [showUserCard, setShowUserCard] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);

    const updateUserListHandler = updateUserList(setData);

    useEffect(() => {
        if (active && !showDeleteDialog) {
            updateUserListHandler();
        }
    }, [active, showDeleteDialog]);

    const handleRowClick = (line: ClientResponse) => {
        setTargetClient(line);
        setShowUserCard(true)
    };

    const handleAdminDeleteUser = async (adminPassword: string) => {
        let resp = null;
        try {
            resp = await adminDeleteUser(targetClient.Email, adminPassword);
            if (resp.Success) {
                setPopupMessage("Success");
            } else {
                setPopupMessage(`Error: ${resp.Error}`);
            }
            setShowDeleteDialog(false);
            setShowUserCard(false);
        } catch (error: any) {
            setPopupMessage(`Error: ${error.message}`);
            console.warn(error.message);
        }
        setShowPopup(true);
    }

    return (
        <>
            <Frame>
                <div className="table-container">
                    <h2>All Users</h2>
                    <GenericTable
                        data={data}
                        columns={['Email', 'Role', 'CreatedDt']}
                        onRowClick={handleRowClick}
                        filterable={true}
                        sortable={true}
                        interpretBigNumbersAs='date'
                    />
                </div>
            </Frame>
            <UserInspectCard
                client={targetClient}
                hidden={!showUserCard}
                onDelete={ ()=>setShowDeleteDialog(true) }
                onClose={()=>setShowUserCard(false)}
            />
            <ConfirmDeleteDialog
                hidden={!showDeleteDialog}
                userName={targetClient.Email}
                onConfirmDelete={handleAdminDeleteUser}
                onCancel={ ()=>setShowDeleteDialog(false) }
            />
        </>
    );
};

function updateUserList(setData: React.Dispatch<React.SetStateAction<ClientResponse[]>>): () => Promise<void> {
    return async () => {
        if (loadingClientList.value == true) {
            return;
        }
        loadingClientList.value = true;
        try {
            await listAllClients()
                .then(value => {
                    if (value != null) {
                        setData(value);
                    }
                });
        } catch (error: any) {
            console.warn(error.message);
        }
        loadingClientList.value = false;
    };
}

export default UserListView;
