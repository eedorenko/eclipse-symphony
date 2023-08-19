import React from 'react'
import { getServerSession } from 'next-auth';
import { options } from '../api/auth/[...nextauth]/options';
import SiteList from '@/components/SiteList';
const getSites = async () => {
  const session = await getServerSession(options);  
  console.log(session?.user?.accessToken);
  const symphonyApi = process.env.SYMPHONY_API;
  const res = await fetch( `${symphonyApi}federation/registry`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${session?.user?.accessToken}`,
    }
  });
  const data = await res.json();
  //map each element and do some transformation
  const sites = data.map((site: any) => {
    return {
      id: site.id,
      name: site.spec.name,
      phone: site.spec.properties['phone'],
      description: site.spec.properties['description'],
    }
  });
  return sites;
}

async function SitesPage() {
  const sites = await getSites();  
  return (
    <div className='cards_view'>
      <SiteList sites={sites}/>      
    </div>
  );
}

export default SitesPage;