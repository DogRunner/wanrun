alter table dogs add constraint dev_dogs_dog_owner_id_fkey foreign key (dog_owner_id) references dog_owners (dog_owner_id);
alter table dogs add constraint dev_dogs_dog_type_id_fkey foreign key (dog_type_id) references dog_type_mst (dog_type_id);

alter table injection_certifications add constraint dev_injection_certifications_dog_id_fkey foreign key (dog_id) references dogs (dog_id);

alter table dogruns add constraint dev_dogruns_dogrun_manager_id_fkey foreign key (dogrun_manager_id) references dogrun_managers (dogrun_manager_id);

alter table regular_business_hours add constraint dev_regular_business_hours_dogrun_id_fkey foreign key (dogrun_id) references dogruns (dogrun_id);

alter table special_business_hours add constraint dev_special_business_hours_dogrun_id_fkey foreign key (dogrun_id) references dogruns (dogrun_id);

alter table dogrun_images add constraint dev_dogrun_images_dogrun_id_fkey foreign key (dogrun_id) references dogruns (dogrun_id);

alter table dogrun_tags add constraint dev_dogrun_tags_dogrun_id_fkey foreign key (dogrun_id) references dogruns (dogrun_id);
alter table dogrun_tags add constraint dev_dogrun_tags_tag_id_fkey foreign key (tag_id) references tag_mst (tag_id);

alter table auth_dog_owners add constraint dev_auth_dog_owners_dog_owner_id_fkey foreign key (dog_owner_id) references dog_owners (dog_owner_id);

alter table dog_owner_credentials add constraint dev_dog_owner_credentials_id_fkey foreign key (auth_dog_owner_id) references auth_dog_owners(auth_dog_owner_id);

alter table dogrun_bookmarks add constraint dev_dogrun_bookmarks_dogrun_id_fkey foreign key (dogrun_id) references dogruns (dogrun_id);
alter table dogrun_bookmarks add constraint dev_dogrun_bookmarks_dog_owner_id_fkey foreign key (dog_owner_id) references dog_owners (dog_owner_id);

alter table s3_file_info add constraint dev_s3_file_info_dog_owners_id_fkey foreign key (dog_owner_id) references dog_owners(dog_owner_id);

-- `organizations`と`dogrun_managers`のリレーション
alter table dogrun_managers add constraint dev_dogrun_managers_organization_id_fkey foreign key (organization_id) references organizations (organization_id);

-- `dogrun_managers`と`auth_dogrun_managers`のリレーション
alter table auth_dogrun_managers add constraint dev_auth_dogrun_managers_dogrun_manager_id_fkey foreign key (dogrun_manager_id) references dogrun_managers (dogrun_manager_id);

-- `auth_dogrun_managers`と`dogrun_manager_credentials`のリレーション
alter table dogrun_manager_credentials add constraint dev_dogrun_manager_credentials_auth_dogrun_manager_id_fkey foreign key (auth_dogrun_manager_id) references auth_dogrun_managers (auth_dogrun_manager_id);
