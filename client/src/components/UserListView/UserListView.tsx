import { useEffect, useState } from 'react';
import { signal } from '@preact/signals-react';

import Frame from '../common/Frame/Frame';
import GenericTable from '../common/GenericTable/GenericTable';
import { listAllClients, type ClientResponse } from '../../api/api';
import UserInspectCard from '../UserInspectCard/UserInspectCard';


const loadingClientList = signal(false);

interface Props {
    active: boolean;
}

function UserListView({active}: Props) {
    const [data, setData] = useState(new Array<ClientResponse>());
    const [client, setClient] = useState({} as ClientResponse); // ClientResponse
    const [showUserCard, setShowUserCard] = useState(false);

    const updateUserListHandler = updateUserList(setData);

    useEffect(() => {
        if (active) {
            updateUserListHandler();
        }
    }, [active]);

    const handleRowClick = (line: ClientResponse) => {
        setClient(line);
        setShowUserCard(true)
    };

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
        <UserInspectCard client={client} hidden={!showUserCard} onDelete={ ()=>console.log("clicked delete") } onClose={()=>setShowUserCard(false)} />
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
